import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "GearboxV2 - Auth",
  description: "GearboxV2 - Auth",
};

export default function AuthLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <section>
      {children}
    </section>
  );
}
