"use client";

import { Clock, RefreshCw, Search } from "lucide-react";
import { useState } from "react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

import { sampleEmails } from "../lib/sample-data";

import { EmailDetail } from "./email-detail";
import { EmailItem } from "./email-item";

import type { Email } from "../lib/types";

export function Inbox() {
  const [emails, setEmails] = useState<Email[]>(sampleEmails);
  const [selectedEmail, setSelectedEmail] = useState<Email | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [refreshing, setRefreshing] = useState(false);

  const unreadCount = emails.filter((email) => !email.read).length;

  const filteredEmails = emails.filter(
    (email) =>
      email.subject.toLowerCase().includes(searchQuery.toLowerCase()) ||
      email.from.toLowerCase().includes(searchQuery.toLowerCase()) ||
      email.content.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const handleRefresh = () => {
    setRefreshing(true);
    // Simulate fetching new emails
    setTimeout(() => {
      setRefreshing(false);
    }, 1000);
  };

  const handleEmailSelect = (email: Email) => {
    // Mark as read when selected
    if (!email.read) {
      const updatedEmails = emails.map((e) =>
        e.id === email.id ? { ...e, read: true } : e
      );
      setEmails(updatedEmails);
    }
    setSelectedEmail(email);
  };

  const handleArchive = (emailId: string) => {
    const updatedEmails = emails.filter((email) => email.id !== emailId);
    setEmails(updatedEmails);
    if (selectedEmail?.id === emailId) {
      setSelectedEmail(null);
    }
  };

  const handleDelete = (emailId: string) => {
    const updatedEmails = emails.filter((email) => email.id !== emailId);
    setEmails(updatedEmails);
    if (selectedEmail?.id === emailId) {
      setSelectedEmail(null);
    }
  };

  return (
    <Card className="h-full">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <div>
          <CardTitle className="text-xl">Inbox</CardTitle>
          <CardDescription>
            {unreadCount > 0 ? (
              <span className="flex items-center gap-2">
                You have <Badge variant="secondary">{unreadCount}</Badge> unread
                messages
              </span>
            ) : (
              "No new messages"
            )}
          </CardDescription>
        </div>
        <Button
          variant="outline"
          size="icon"
          onClick={handleRefresh}
          disabled={refreshing}
        >
          <RefreshCw
            className={`h-4 w-4 ${refreshing ? "animate-spin" : ""}`}
          />
          <span className="sr-only">Refresh inbox</span>
        </Button>
      </CardHeader>
      <CardContent>
        <div className="mb-4">
          <div className="relative">
            <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search emails..."
              className="pl-8"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>
        </div>

        <Tabs defaultValue="all" className="space-y-4">
          <TabsList>
            <TabsTrigger value="all">All</TabsTrigger>
            <TabsTrigger value="unread">
              Unread
              {unreadCount > 0 && (
                <Badge variant="secondary" className="ml-2">
                  {unreadCount}
                </Badge>
              )}
            </TabsTrigger>
          </TabsList>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
            <div className="lg:col-span-1 space-y-2 overflow-auto max-h-[600px] pr-2">
              <TabsContent value="all" className="m-0 space-y-2">
                {filteredEmails.length > 0 ? (
                  filteredEmails.map((email) => (
                    <EmailItem
                      key={email.id}
                      email={email}
                      isSelected={selectedEmail?.id === email.id}
                      onSelect={() => handleEmailSelect(email)}
                      onArchive={() => handleArchive(email.id)}
                      onDelete={() => handleDelete(email.id)}
                    />
                  ))
                ) : (
                  <div className="flex flex-col items-center justify-center py-12 text-center text-muted-foreground">
                    {/* <Inbox className="h-12 w-12 mb-4" /> */}
                    <h3 className="text-lg font-medium">No emails found</h3>
                    <p className="text-sm">
                      {searchQuery
                        ? "Try a different search term"
                        : "Your inbox is empty"}
                    </p>
                  </div>
                )}
              </TabsContent>

              <TabsContent value="unread" className="m-0 space-y-2">
                {filteredEmails.filter((e) => !e.read).length > 0 ? (
                  filteredEmails
                    .filter((email) => !email.read)
                    .map((email) => (
                      <EmailItem
                        key={email.id}
                        email={email}
                        isSelected={selectedEmail?.id === email.id}
                        onSelect={() => handleEmailSelect(email)}
                        onArchive={() => handleArchive(email.id)}
                        onDelete={() => handleDelete(email.id)}
                      />
                    ))
                ) : (
                  <div className="flex flex-col items-center justify-center py-12 text-center text-muted-foreground">
                    {/* <Inbox className="h-12 w-12 mb-4" /> */}
                    <h3 className="text-lg font-medium">No unread emails</h3>
                    <p className="text-sm">
                      {searchQuery
                        ? "Try a different search term"
                        : "You're all caught up!"}
                    </p>
                  </div>
                )}
              </TabsContent>
            </div>

            <div className="lg:col-span-2 border rounded-lg p-4 bg-card min-h-[400px]">
              {selectedEmail ? (
                <EmailDetail
                  email={selectedEmail}
                  onArchive={() => handleArchive(selectedEmail.id)}
                  onDelete={() => handleDelete(selectedEmail.id)}
                />
              ) : (
                <div className="flex flex-col items-center justify-center h-full py-12 text-center text-muted-foreground">
                  <Clock className="h-12 w-12 mb-4" />
                  <h3 className="text-lg font-medium">No email selected</h3>
                  <p className="text-sm">
                    Select an email to view its contents
                  </p>
                </div>
              )}
            </div>
          </div>
        </Tabs>
      </CardContent>
    </Card>
  );
}
