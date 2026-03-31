import React from "react"
import Link from "next/link"
import {
  Zap,
  Users,
  Code2,
  Star,
  Shield,
  BookOpen,
  Newspaper,
  Compass,
  MessageCircle,
  ArrowRight,
  TrendingUp,
  Award,
  HeartHandshake,
} from "lucide-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

// ── Data ──────────────────────────────────────────────────────────────────────

const stats = [
  { value: "10,000+", label: "Happy Customers", icon: Users },
  { value: "500+",    label: "Expert Reviews",   icon: Star },
  { value: "50+",     label: "Tech Authors",      icon: Award },
  { value: "99%",     label: "Satisfaction Rate", icon: HeartHandshake },
]

const features = [
  {
    icon: BookOpen,
    title: "In-Depth Reviews",
    description:
      "Thorough, unbiased reviews written by hands-on tech experts so you can shop with confidence.",
  },
  {
    icon: Newspaper,
    title: "Tech News",
    description:
      "Stay ahead of the curve with breaking industry news, product launches, and trend analysis.",
  },
  {
    icon: Compass,
    title: "Expert Guides",
    description:
      "Step-by-step guides and buying advice tailored for every skill level — from beginner to pro.",
  },
  {
    icon: MessageCircle,
    title: "Community",
    description:
      "Join thousands of tech enthusiasts to share tips, ask questions, and celebrate great gear.",
  },
  {
    icon: Shield,
    title: "Trusted & Transparent",
    description:
      "No sponsored rankings, no hidden agendas — just honest assessments you can rely on.",
  },
  {
    icon: Zap,
    title: "Fast & Accessible",
    description:
      "Optimised for speed on every device so you get the information you need, instantly.",
  },
]

const values = [
  {
    icon: Shield,
    title: "Integrity",
    description: "We never compromise our editorial independence. Every opinion is genuinely ours.",
  },
  {
    icon: Star,
    title: "Excellence",
    description: "From testing methodology to writing, we hold ourselves to the highest standard.",
  },
  {
    icon: HeartHandshake,
    title: "Community First",
    description: "Our readers shape our roadmap. Your feedback drives every decision we make.",
  },
]

// ── Page ──────────────────────────────────────────────────────────────────────

export default function AboutPage() {
  return (
    <div className="min-h-screen bg-background">

      {/* ── Hero ──────────────────────────────────────────────────────────── */}
      <section className="relative overflow-hidden py-24 lg:py-32">
        {/* Decorative background blob */}
        <div
          aria-hidden
          className="pointer-events-none absolute inset-0 -z-10 flex items-center justify-center"
        >
          <div className="h-[600px] w-[600px] rounded-full bg-primary/5 blur-3xl" />
        </div>

        <div className="container mx-auto px-4 text-center">
          <Badge variant="secondary" className="mb-6">
            <TrendingUp className="mr-2 h-4 w-4" />
            Your Most Trusted Tech Companion
          </Badge>

          <h1 className="mb-6 text-4xl font-bold tracking-tight text-foreground md:text-6xl lg:text-7xl">
            About{" "}
            <span className="bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
              GearBox
            </span>
          </h1>

          <p className="mx-auto mb-10 max-w-2xl text-lg text-muted-foreground leading-relaxed">
            We're a team of passionate tech lovers on a mission to make technology accessible,
            understandable, and exciting for everyone — from first-time buyers to seasoned enthusiasts.
          </p>

          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button size="lg" asChild>
              <Link href="/shop">
                Explore Products
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
            <Button size="lg" variant="outline" asChild>
              <Link href="/blogs">Read Our Reviews</Link>
            </Button>
          </div>
        </div>
      </section>

      {/* ── Stats ─────────────────────────────────────────────────────────── */}
      <section className="py-16 bg-muted/30">
        <div className="container mx-auto px-4">
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
            {stats.map(({ value, label, icon: Icon }) => (
              <Card
                key={label}
                className="group text-center hover:shadow-lg transition-all duration-300 hover:-translate-y-1"
              >
                <CardContent className="pt-8 pb-6">
                  <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-primary/10 group-hover:bg-primary/20 transition-colors duration-300">
                    <Icon className="h-7 w-7 text-primary" />
                  </div>
                  <div className="text-4xl font-bold text-primary mb-1">{value}</div>
                  <div className="text-sm text-muted-foreground font-medium">{label}</div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* ── Mission ───────────────────────────────────────────────────────── */}
      <section className="py-20">
        <div className="container mx-auto px-4">
          <div className="max-w-3xl mx-auto text-center">
            <Badge variant="secondary" className="mb-4">Our Mission</Badge>
            <h2 className="text-3xl font-bold text-foreground mb-6 md:text-4xl">
              Making Technology Work{" "}
              <span className="bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
                for You
              </span>
            </h2>
            <p className="text-muted-foreground leading-relaxed text-lg">
              At GearBox, we believe great technology should be within everyone's reach. That means
              cutting through the marketing noise, providing genuinely useful guidance, and building a
              community where curiosity and knowledge are celebrated. Whether you're searching for your
              first laptop or the latest flagship phone, we have the insight to help you choose wisely.
            </p>
          </div>
        </div>
      </section>

      {/* ── What We Do ────────────────────────────────────────────────────── */}
      <section className="py-16 bg-muted/30">
        <div className="container mx-auto px-4">
          <div className="mb-12 text-center">
            <Badge variant="secondary" className="mb-4">What We Do</Badge>
            <h2 className="text-3xl font-bold text-foreground mb-4 md:text-4xl">
              Everything Tech, All in One Place
            </h2>
            <p className="text-muted-foreground max-w-2xl mx-auto">
              From the latest reviews to deep-dive guides, we've got every angle covered.
            </p>
          </div>

          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {features.map(({ icon: Icon, title, description }) => (
              <Card
                key={title}
                className="group hover:shadow-lg transition-all duration-300 hover:-translate-y-1"
              >
                <CardHeader>
                  <div className="mb-3 flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 group-hover:bg-primary/20 transition-colors duration-300">
                    <Icon className="h-6 w-6 text-primary" />
                  </div>
                  <CardTitle className="text-lg">{title}</CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted-foreground leading-relaxed">{description}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* ── Our Values ────────────────────────────────────────────────────── */}
      <section className="py-20">
        <div className="container mx-auto px-4">
          <div className="mb-12 text-center">
            <Badge variant="secondary" className="mb-4">Our Values</Badge>
            <h2 className="text-3xl font-bold text-foreground mb-4 md:text-4xl">
              Principles We Live By
            </h2>
            <p className="text-muted-foreground max-w-2xl mx-auto">
              Everything we do flows from a simple set of beliefs about what good tech journalism looks like.
            </p>
          </div>

          <div className="grid gap-8 md:grid-cols-3 max-w-4xl mx-auto">
            {values.map(({ icon: Icon, title, description }) => (
              <div key={title} className="flex flex-col items-center text-center gap-4">
                <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-primary/10">
                  <Icon className="h-8 w-8 text-primary" />
                </div>
                <div>
                  <h3 className="font-semibold text-foreground mb-2 text-lg">{title}</h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">{description}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ── Solo Builder ──────────────────────────────────────────────────── */}
      <section className="py-16 bg-muted/30">
        <div className="container mx-auto px-4">
          <div className="max-w-3xl mx-auto">
            <Card className="overflow-hidden">
              <CardContent className="pt-10 pb-10">
                <div className="flex flex-col sm:flex-row items-center sm:items-start gap-8">
                  {/* Avatar */}
                  <div className="shrink-0 flex h-24 w-24 items-center justify-center rounded-2xl bg-primary/10">
                    <Code2 className="h-12 w-12 text-primary" />
                  </div>

                  {/* Text */}
                  <div className="text-center sm:text-left">
                    <Badge variant="secondary" className="mb-3">Solo Builder</Badge>
                    <h2 className="text-2xl font-bold text-foreground mb-3">
                      Built by One Passionate Developer
                    </h2>
                    <p className="text-muted-foreground leading-relaxed">
                      GearBox is a solo project — designed, built, and maintained by a single developer
                      who loves tech as much as you do. Every feature, review workflow, and line of code
                      exists because it solves a real problem. No committee, no compromises — just a
                      genuine passion for great technology and honest information.
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* ── CTA ───────────────────────────────────────────────────────────── */}
      <section className="py-20">
        <div className="container mx-auto px-4">
          <Card className="max-w-3xl mx-auto overflow-hidden">
            <div className="relative">
              {/* Subtle gradient accent */}
              <div
                aria-hidden
                className="pointer-events-none absolute inset-0 bg-gradient-to-br from-primary/5 to-transparent"
              />
              <CardContent className="relative pt-12 pb-12 text-center">
                <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-primary/10">
                  <Zap className="h-8 w-8 text-primary" />
                </div>
                <h2 className="text-3xl font-bold text-foreground mb-4">Ready to Geek Out?</h2>
                <p className="text-muted-foreground max-w-md mx-auto mb-8 leading-relaxed">
                  Browse our curated product catalogue, dive into expert reviews, or join the
                  conversation with thousands of fellow tech fans.
                </p>
                <div className="flex flex-col sm:flex-row gap-4 justify-center">
                  <Button size="lg" asChild>
                    <Link href="/shop">
                      Shop Now
                      <ArrowRight className="ml-2 h-4 w-4" />
                    </Link>
                  </Button>
                  <Button size="lg" variant="outline" asChild>
                    <Link href="/contact">Get in Touch</Link>
                  </Button>
                </div>
              </CardContent>
            </div>
          </Card>
        </div>
      </section>

    </div>
  )
}
