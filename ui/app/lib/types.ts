import type { LucideIcon } from "lucide-react"

export interface Email {
  id: string
  from: string
  subject: string
  content: string
  timestamp: Date
  read: boolean
  attachments?: {
    name: string
    size: string
    icon: LucideIcon
  }[]
}

