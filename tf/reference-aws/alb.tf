# ALB + target group pointing at the EC2 instance.
resource "aws_lb" "ref" {
  name               = "domovyk-ref-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.web.id]
  subnets            = aws_subnet.public[*].id
  enable_deletion_protection = false

  tags = merge(local.common_tags, {
    Name = "domovyk-ref-alb"
  })
}

resource "aws_lb_target_group" "ref" {
  name     = "domovyk-ref-tg"
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.ref.id
  health_check {
    path = "/"
  }
  tags = local.common_tags
}

resource "aws_lb_listener" "ref" {
  load_balancer_arn = aws_lb.ref.id
  port              = "80"
  protocol          = "HTTP"
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.ref.id
  }
}

resource "aws_lb_target_group_attachment" "ref" {
  target_group_arn = aws_lb_target_group.ref.id
  target_id        = aws_instance.web.id
  port             = 80
}
