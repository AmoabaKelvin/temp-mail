import { File, FileText, Image } from "lucide-react"
import type { Email } from "./types"

export const sampleEmails: Email[] = [
  {
    id: "1",
    from: "notifications@github.com",
    subject: "New pull request for your repository",
    content:
      "A new pull request has been opened in your repository by user123.\n\nPull Request #42: Fix navigation bug in mobile view\n\nPlease review the changes and provide feedback.",
    timestamp: new Date(Date.now() - 1000 * 60 * 15), // 15 minutes ago
    read: false,
  },
  {
    id: "2",
    from: "support@vercel.com",
    subject: "Your deployment is complete",
    content:
      "Your project has been successfully deployed to production.\n\nDeployment URL: https://your-project.vercel.app\n\nView your deployment dashboard for more details.",
    timestamp: new Date(Date.now() - 1000 * 60 * 45), // 45 minutes ago
    read: false,
  },
  {
    id: "3",
    from: "newsletter@reactjs.org",
    subject: "React Newsletter - March 2025",
    content:
      "Here are the latest updates from the React ecosystem:\n\n- React 20 released with new features\n- Community spotlight: Building accessible components\n- Upcoming React conferences and events\n\nRead more on our blog.",
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 3), // 3 hours ago
    read: true,
  },
  {
    id: "4",
    from: "billing@stripe.com",
    subject: "Your monthly receipt",
    content:
      "Thank you for your subscription. Here is your receipt for March 2025.\n\nAmount: $15.00\nDate: March 15, 2025\nPayment method: •••• 4242\n\nView your billing dashboard for more details.",
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 8), // 8 hours ago
    read: true,
    attachments: [
      {
        name: "receipt-march-2025.pdf",
        size: "156 KB",
        icon: FileText,
      },
    ],
  },
  {
    id: "5",
    from: "team@tailwindcss.com",
    subject: "Tailwind CSS v5.0 is here!",
    content:
      "We're excited to announce the release of Tailwind CSS v5.0!\n\nHighlights:\n- Improved performance\n- New utility classes\n- Better dark mode support\n- Enhanced documentation\n\nCheck out our blog post for all the details.",
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24), // 1 day ago
    read: true,
  },
  {
    id: "6",
    from: "no-reply@aws.amazon.com",
    subject: "Your AWS usage report",
    content:
      "Here is your AWS usage report for the current billing period.\n\nCurrent charges: $42.50\nEstimated month-end: $78.20\n\nServices:\n- EC2: $28.15\n- S3: $8.75\n- Lambda: $5.60\n\nLog in to your AWS console for more details.",
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24 * 2), // 2 days ago
    read: true,
    attachments: [
      {
        name: "aws-report-march.csv",
        size: "24 KB",
        icon: File,
      },
    ],
  },
  {
    id: "7",
    from: "design@figma.com",
    subject: "Your team has new comments",
    content:
      "Your team has added new comments to the design file 'Website Redesign'.\n\nJane: Can we adjust the spacing on the hero section?\nMike: The color contrast needs improvement for accessibility.\n\nView all comments in Figma.",
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24 * 3), // 3 days ago
    read: true,
    attachments: [
      {
        name: "design-preview.png",
        size: "1.2 MB",
        icon: Image,
      },
    ],
  },
]

