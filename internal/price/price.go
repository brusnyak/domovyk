// Package price compares one exported month's consumption with supplied EUR rates.
// It does not infer equivalent resource sizes or convert incompatible units.
package price

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"gopkg.in/yaml.v3"
)

type Price struct {
	Provider  string   `yaml:"provider"`
	SKU       string   `yaml:"sku"`
	Category  string   `yaml:"category"`
	Unit      string   `yaml:"unit"`
	PriceEUR  *float64 `yaml:"price_eur"`
	SourceURL string   `yaml:"source_url"`
	CheckedOn string   `yaml:"checked_on"`
	Status    string   `yaml:"status"`
}

type Result struct {
	AWSService       string   `json:"aws_service"`
	AWSMonthly       float64  `json:"aws_monthly_eur_or_billing_currency"`
	Currency         string   `json:"currency"`
	Provider         string   `json:"provider"`
	SKU              string   `json:"sku"`
	EstimatedMonthly *float64 `json:"estimated_monthly"`
	Delta            *float64 `json:"delta"` // AWS minus target; positive means savings.
	CheckedOn        string   `json:"checked_on"`
	Status           string   `json:"status"`
	SourceURL        string   `json:"source_url"`
	Reason           string   `json:"reason,omitempty"`
}

func number(s string) (float64, error) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(n) || math.IsInf(n, 0) {
		return 0, fmt.Errorf("invalid finite number %q", s)
	}
	return n, nil
}

// ParseFX accepts an explicit billing-currency-to-EUR multiplier.
func ParseFX(s string) (map[string]float64, error) {
	fx := map[string]float64{}
	if s == "" {
		return fx, nil
	}
	parts := strings.Split(s, "=")
	if len(parts) != 2 || len(parts[0]) != 3 || parts[0] == "EUR" {
		return nil, fmt.Errorf("fx must be CURRENCY=positive_EUR_rate, e.g. USD=0.92")
	}
	for _, c := range parts[0] {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid fx currency")
		}
	}
	rate, err := number(parts[1])
	if err != nil || rate <= 0 {
		return nil, fmt.Errorf("fx requires a positive finite EUR rate")
	}
	fx[parts[0]] = rate
	return fx, nil
}

func validCategory(c string) bool {
	switch c {
	case "compute", "block_storage", "object_storage", "managed_db", "egress":
		return true
	}
	return false
}

func LoadPrices(path string) ([]Price, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var file struct {
		Rows []Price `yaml:"rows"`
	}
	dec := yaml.NewDecoder(strings.NewReader(string(data)))
	dec.KnownFields(true)
	if err := dec.Decode(&file); err != nil {
		return nil, fmt.Errorf("prices: %w", err)
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("prices: expected one YAML document")
	}
	for i, p := range file.Rows {
		_, dateErr := time.Parse("2006-01-02", p.CheckedOn)
		u, urlErr := url.Parse(p.SourceURL)
		validURL := urlErr == nil && u.Host != "" && (u.Scheme == "https" || u.Scheme == "http")
		if (p.Provider != "hetzner" && p.Provider != "ovhcloud") || strings.TrimSpace(p.SKU) == "" || !validCategory(p.Category) || strings.TrimSpace(p.Unit) == "" || (p.PriceEUR == nil && p.Status != "unverified") || (p.PriceEUR != nil && (*p.PriceEUR < 0 || math.IsNaN(*p.PriceEUR) || math.IsInf(*p.PriceEUR, 0))) || dateErr != nil || !validURL || (p.Status != "verified" && p.Status != "unverified") {
			return nil, fmt.Errorf("prices row %d: invalid provider, SKU, category, unit, rate or provenance", i+1)
		}
	}
	return file.Rows, nil
}

// The seed map has prose targets, so only known, explicitly named target
// categories are supported. Software rewrites and networking targets are skipped.
func loadTargets(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var file struct {
		Rows []struct {
			ID      string `yaml:"id"`
			Service string `yaml:"aws_service"`
			Target  string `yaml:"typical_eu_target"`
		} `yaml:"rows"`
	}
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, err
	}
	categories := map[string]string{"aws-ec2-instance": "compute", "aws-ebs-volume": "block_storage", "aws-s3-bucket": "object_storage", "aws-rds-postgresql": "managed_db"}
	aliases := map[string][]string{
		"aws-ec2-instance":   {"Amazon EC2", "AmazonEC2", "Amazon Elastic Compute Cloud"},
		"aws-ebs-volume":     {"Amazon EBS", "AmazonEBS", "Amazon Elastic Block Store"},
		"aws-s3-bucket":      {"Amazon S3", "AmazonS3", "Amazon Simple Storage Service"},
		"aws-rds-postgresql": {"Amazon RDS", "AmazonRDS", "Amazon Relational Database Service"},
	}
	targets := map[string]string{}
	for _, r := range file.Rows {
		c := categories[r.ID]
		if c == "" || strings.TrimSpace(r.Target) == "" {
			continue
		}
		targets[strings.ToLower(r.Service)] = c
		for _, alias := range aliases[r.ID] {
			targets[strings.ToLower(alias)] = c
		}
	}
	return targets, nil
}

type group struct {
	cost         float64
	quantities   map[string]float64
	mixedCompute bool
}

func unit(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "h", "hr", "hrs", "hour", "hours":
		return "hour"
	default:
		return strings.ToLower(strings.TrimSpace(s))
	}
}

func readFocus(path string, fx map[string]float64) (map[string]*group, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("FOCUS header: %w", err)
	}
	columns := map[string]int{}
	for i, name := range header {
		name = strings.TrimPrefix(name, "\ufeff")
		if _, exists := columns[name]; exists {
			return nil, fmt.Errorf("duplicate FOCUS column %s", name)
		}
		columns[name] = i
	}
	for _, name := range strings.Split("ServiceName,ServiceCategory,ResourceId,ResourceType,ChargeCategory,BilledCost,EffectiveCost,BillingCurrency,ConsumedQuantity,ConsumedUnit,ChargePeriodStart", ",") {
		if _, ok := columns[name]; !ok {
			return nil, fmt.Errorf("missing FOCUS column %s", name)
		}
	}
	groups := map[string]*group{}
	month := ""
	for line := 2; ; line++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("FOCUS row %d: %w", line, err)
		}
		get := func(name string) string { return strings.TrimSpace(record[columns[name]]) }
		service := get("ServiceName")
		if service == "" {
			return nil, fmt.Errorf("FOCUS row %d: empty ServiceName", line)
		}
		cost, err := number(get("BilledCost"))
		if err != nil {
			return nil, fmt.Errorf("FOCUS row %d BilledCost: %w", line, err)
		}
		if _, err := number(get("EffectiveCost")); err != nil {
			return nil, fmt.Errorf("FOCUS row %d EffectiveCost: %w", line, err)
		}
		currency, rate := get("BillingCurrency"), 1.0
		if currency != "EUR" {
			var ok bool
			rate, ok = fx[currency]
			if !ok || rate <= 0 || math.IsNaN(rate) || math.IsInf(rate, 0) {
				return nil, fmt.Errorf("FOCUS row %d: currency %q requires explicit --fx %s=EUR_rate", line, currency, currency)
			}
		}
		date, err := time.Parse(time.RFC3339, get("ChargePeriodStart"))
		if err != nil {
			return nil, fmt.Errorf("FOCUS row %d: invalid ChargePeriodStart: %w", line, err)
		}
		m := date.Format("2006-01")
		if month != "" && m != month {
			return nil, fmt.Errorf("FOCUS spans multiple months; supply one calendar month")
		}
		month = m
		g := groups[service]
		if g == nil {
			g = &group{quantities: map[string]float64{}}
			groups[service] = g
		}
		g.cost += cost * rate
		if math.IsInf(g.cost, 0) {
			return nil, fmt.Errorf("FOCUS row %d: cost overflow", line)
		}
		// Credits/purchases affect billed cost but do not add consumption.
		if get("ChargeCategory") != "Usage" {
			continue
		}
		q, err := number(get("ConsumedQuantity"))
		if err != nil || q < 0 {
			return nil, fmt.Errorf("FOCUS row %d: invalid ConsumedQuantity", line)
		}
		if q > 0 {
			g.quantities[unit(get("ConsumedUnit"))] += q
		}
		// EC2 exports can contain EBS and transfer. Do not price those as compute.
		rt := strings.ToLower(get("ResourceType"))
		if strings.Contains(rt, "volume") || strings.Contains(rt, "snapshot") || strings.Contains(rt, "transfer") {
			g.mixedCompute = true
		}
	}
	return groups, nil
}

// Compare sums BilledCost, including credits, and multiplies exported consumption
// by compatible unit rates. no_price uses table provenance and an explicit
// not_available check date rather than inventing a price verification date.
func Compare(focusPath, pricesPath, mappingPath string, fx map[string]float64) ([]Result, error) {
	prices, err := LoadPrices(pricesPath)
	if err != nil {
		return nil, err
	}
	targets, err := loadTargets(mappingPath)
	if err != nil {
		return nil, err
	}
	groups, err := readFocus(focusPath, fx)
	if err != nil {
		return nil, err
	}
	abs, err := filepath.Abs(pricesPath)
	if err != nil {
		return nil, err
	}
	provenance := (&url.URL{Scheme: "file", Path: abs}).String()
	services := make([]string, 0, len(groups))
	for s := range groups {
		services = append(services, s)
	}
	sort.Strings(services)
	results := []Result{}
	for _, service := range services {
		category := targets[strings.ToLower(service)]
		if category == "" {
			continue
		}
		g, start := groups[service], len(results)
		for _, p := range prices {
			if p.PriceEUR == nil || p.Category != category || len(g.quantities) != 1 || (category == "compute" && g.mixedCompute) {
				continue
			}
			q, ok := g.quantities[unit(p.Unit)]
			if !ok {
				continue
			}
			estimate := q * *p.PriceEUR
			delta := g.cost - estimate
			if math.IsNaN(estimate) || math.IsInf(estimate, 0) || math.IsInf(delta, 0) {
				return nil, fmt.Errorf("estimate overflow for %s", service)
			}
			results = append(results, Result{AWSService: service, AWSMonthly: g.cost, Currency: "EUR", Provider: p.Provider, SKU: p.SKU, EstimatedMonthly: &estimate, Delta: &delta, CheckedOn: p.CheckedOn, Status: p.Status, SourceURL: p.SourceURL})
		}
		if len(results) == start {
			results = append(results, Result{AWSService: service, AWSMonthly: g.cost, Currency: "EUR", CheckedOn: "not_available", Status: "no_price", SourceURL: provenance, Reason: "no compatible price for mapped category " + category + " and exported consumption units"})
		}
	}
	return results, nil
}

// WriteResults atomically replaces the output with mode 0600, even if it existed.
func WriteResults(path string, results []Result) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".domovyk-price-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(append(data, '\n')); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}

func PrintTable(w io.Writer, results []Result) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "AWS SERVICE\tAWS EUR\tPROVIDER\tSKU\tESTIMATED EUR\tDELTA EUR\tCHECKED ON\tSTATUS\tSOURCE URL")
	for _, r := range results {
		estimate, delta := "-", "-"
		if r.EstimatedMonthly != nil {
			estimate = fmt.Sprintf("%.2f", *r.EstimatedMonthly)
		}
		if r.Delta != nil {
			delta = fmt.Sprintf("%.2f", *r.Delta)
		}
		fmt.Fprintf(tw, "%s\t%.2f\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", r.AWSService, r.AWSMonthly, r.Provider, r.SKU, estimate, delta, r.CheckedOn, strings.ToUpper(r.Status), r.SourceURL)
	}
	return tw.Flush()
}
