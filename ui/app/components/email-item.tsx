"use client";

import DOMPurify from "dompurify";
import { Archive, Circle, MailOpen, MailX, Trash2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

import { formatRelativeTime } from "../lib/date-utils";

import type { Email } from "../lib/types";

interface EmailItemProps {
  email: Email;
  isSelected: boolean;
  onSelect: () => void;
  onArchive: () => void;
  onDelete: () => void;
  onToggleRead: () => void;
}

export function EmailItem({
  email,
  isSelected,
  onSelect,
  onArchive,
  onDelete,
  onToggleRead,
}: EmailItemProps) {
  return (
    <div
      className={cn(
        "flex flex-col rounded-lg border p-3 transition-colors hover:bg-accent cursor-pointer",
        isSelected && "bg-accent",
        !email.read && "border-l-4 border-l-primary",
      )}
      onClick={onSelect}
    >
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-2">
          {!email.read && (
            <Circle className="h-2 w-2 fill-primary stroke-primary" />
          )}
          <div className="font-medium truncate max-w-[150px]">{email.from}</div>
        </div>
        <div className="text-xs text-muted-foreground">
          {formatRelativeTime(email.timestamp)}
        </div>
      </div>
      <div className="mt-1">
        <div
          className={cn("font-medium truncate", email.read && "font-normal")}
        >
          {email.subject}
        </div>
        <div
          className="mt-1 text-xs text-muted-foreground line-clamp-1"
          dangerouslySetInnerHTML={{
            __html: DOMPurify.sanitize(email.content),
          }}
        />
      </div>
      <div className="mt-2 flex justify-end gap-1">
        <Button
          variant="ghost"
          size="icon"
          className="h-7 w-7"
          onClick={(e) => {
            e.stopPropagation();
            onToggleRead();
          }}
          title={email.read ? "Mark as unread" : "Mark as read"}
        >
          {email.read ? (
            <MailX className="h-4 w-4" />
          ) : (
            <MailOpen className="h-4 w-4" />
          )}
          <span className="sr-only">
            {email.read ? "Mark as unread" : "Mark as read"}
          </span>
        </Button>
        <Button
          variant="ghost"
          size="icon"
          className="h-7 w-7"
          onClick={(e) => {
            e.stopPropagation();
            onArchive();
          }}
        >
          <Archive className="h-4 w-4" />
          <span className="sr-only">Archive</span>
        </Button>
        <Button
          variant="ghost"
          size="icon"
          className="h-7 w-7"
          onClick={(e) => {
            e.stopPropagation();
            onDelete();
          }}
        >
          <Trash2 className="h-4 w-4" />
          <span className="sr-only">Delete</span>
        </Button>
      </div>
    </div>
  );
}
