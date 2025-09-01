import { LayoutWrapper } from "@/components/layout-wrapper"

export default function StaticLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <LayoutWrapper>
      {children}
    </LayoutWrapper>
  )
}