"use client";

import { Button } from "@/components/ui/button";
import DOMPurify from "dompurify";
import { Archive, Calendar, MailOpen, MailX, Trash2, User } from "lucide-react";

import { formatDateTime } from "../lib/date-utils";

import type { Email } from "../lib/types";

interface EmailDetailProps {
	email: Email;
	onArchive: () => void;
	onDelete: () => void;
	onToggleRead: () => void;
}

export function EmailDetail({
	email,
	onArchive,
	onDelete,
	onToggleRead,
}: EmailDetailProps) {
	return (
		<div className="space-y-4 h-full">
			<div className="flex justify-between items-start">
				<h3 className="text-xl font-bold">{email.subject}</h3>
				<div className="flex gap-1">
					<Button
						variant="outline"
						size="icon"
						onClick={onToggleRead}
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
					<Button variant="outline" size="icon" onClick={onArchive}>
						<Archive className="h-4 w-4" />
						<span className="sr-only">Archive</span>
					</Button>
					<Button variant="outline" size="icon" onClick={onDelete}>
						<Trash2 className="h-4 w-4" />
						<span className="sr-only">Delete</span>
					</Button>
				</div>
			</div>

			<div className="flex flex-col space-y-1 text-sm bg-muted/50 p-3 rounded-md">
				<div className="flex items-center gap-2">
					<User className="h-4 w-4 text-muted-foreground" />
					<span className="font-medium">From:</span> {email.from}
				</div>
				<div className="flex items-center gap-2">
					<Calendar className="h-4 w-4 text-muted-foreground" />
					<span className="font-medium">Date:</span>{" "}
					{formatDateTime(email.timestamp)}
				</div>
			</div>

			<div className="border-t pt-4 flex-1 min-h-[300px]">
				<div
					className="prose prose-sm max-w-none"
					dangerouslySetInnerHTML={{
						__html: DOMPurify.sanitize(email.content),
					}}
				/>
				{email.attachments && email.attachments.length > 0 && (
					<div className="mt-6 bg-muted/30 p-4 rounded-md">
						<h4 className="text-sm font-medium mb-2">Attachments</h4>
						<div className="grid grid-cols-2 gap-2">
							{email.attachments.map((attachment, index) => (
								<div
									key={index}
									className="flex items-center gap-2 rounded-md border p-2 text-sm bg-background hover:bg-accent transition-colors"
								>
									<div className="rounded bg-primary/10 p-1">
										<attachment.icon className="h-4 w-4 text-primary" />
									</div>
									<div className="flex-1 truncate">{attachment.name}</div>
								</div>
							))}
						</div>
					</div>
				)}
			</div>
		</div>
	);
}
