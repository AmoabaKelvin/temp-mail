import { Space_Grotesk } from "next/font/google";

import { EmailHeader } from "@/app/components/email-header";
import { Inbox } from "@/app/components/inbox";

const spaceGrotesk = Space_Grotesk({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
});

export default function Home() {
  return (
    <main className={`h-dvh bg-white ${spaceGrotesk.className}`}>
      <div className="container mx-auto px-4 py-8 flex flex-col gap-8">
        {/* Updated Email Header */}
        <EmailHeader />

        {/* Inbox remains unchanged */}
        <div className="flex-1">
          <Inbox />
        </div>
      </div>
    </main>
  );
}
