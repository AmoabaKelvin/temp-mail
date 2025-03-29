"use client";

import { Clock, RefreshCw, Search } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { toast } from "@/components/ui/use-toast";

import {
  deleteMessage,
  EmailMessage,
  getMessages,
  updateMessageReadStatus,
} from "../lib/api-client";

import { EmailDetail } from "./email-detail";
import { EmailItem } from "./email-item";

import type { Email } from "../lib/types";

const POLLING_INTERVAL_MS = 30000;
const LOCAL_STORAGE_EMAIL_KEY = "currentEmail";

// Convert API EmailMessage to UI Email format
const convertToEmail = (apiMessage: EmailMessage): Email => {
  return {
    id: apiMessage.id.toString(),
    from: apiMessage.from_address,
    subject: apiMessage.subject || "No Subject",
    content: apiMessage.body,
    timestamp: new Date(apiMessage.received_at),
    read: apiMessage.read_at !== null,
  };
};

export function Inbox() {
  const [emails, setEmails] = useState<Email[]>([]);
  const [selectedEmail, setSelectedEmail] = useState<Email | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [refreshing, setRefreshing] = useState(false);
  const [currentEmail, setCurrentEmail] = useState<string | null>(null);
  const [pollingInterval, setPollingInterval] = useState<NodeJS.Timeout | null>(
    null
  );

  const unreadCount = emails.filter((email) => !email.read).length;

  const filteredEmails = useMemo(() => {
    return emails.filter(
      (email) =>
        email.subject.toLowerCase().includes(searchQuery.toLowerCase()) ||
        email.from.toLowerCase().includes(searchQuery.toLowerCase()) ||
        email.content.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }, [emails, searchQuery]);

  // Load current email from localStorage on initial load
  useEffect(() => {
    const storedEmail = localStorage.getItem(LOCAL_STORAGE_EMAIL_KEY);
    if (storedEmail) {
      setCurrentEmail(storedEmail);
      fetchMessages(storedEmail);
    }
  }, []);

  // Fetch messages when current email changes
  const fetchMessages = async (email: string) => {
    if (!email) return;

    setRefreshing(true);
    const response = await getMessages(email);
    setRefreshing(false);

    if (response.error) {
      toast({
        title: "Error loading messages",
        description: response.error.message,
        variant: "destructive",
      });
      return;
    }

    if (response.data) {
      // Convert API messages to UI Email format
      const convertedEmails = response.data.map(convertToEmail);
      setEmails(convertedEmails);
    }
  };

  // Set up polling for new emails every 30 seconds
  useEffect(() => {
    if (currentEmail) {
      // Clear any existing interval
      if (pollingInterval) {
        clearInterval(pollingInterval);
      }

      fetchMessages(currentEmail);

      // Set up new polling interval (every 30 seconds)
      const interval = setInterval(() => {
        setRefreshing(true);
        fetchMessages(currentEmail);
      }, POLLING_INTERVAL_MS);

      setPollingInterval(interval);

      // Clean up interval on unmount or when currentEmail changes
      return () => {
        if (interval) clearInterval(interval);
      };
    }
  }, [currentEmail]);

  const handleRefresh = () => {
    if (currentEmail) {
      fetchMessages(currentEmail);
    } else {
      toast({
        title: "No email selected",
        description: "Generate an email address first.",
        variant: "destructive",
      });
    }
  };

  const handleEmailSelect = async (email: Email) => {
    // Mark as read when selected
    if (!email.read) {
      // Update UI optimistically
      const updatedEmails = emails.map((e) =>
        e.id === email.id ? { ...e, read: true } : e
      );
      setEmails(updatedEmails);

      // Call API to update read status
      const response = await updateMessageReadStatus(parseInt(email.id));
      if (response.error) {
        // Revert UI change if API call fails
        const revertedEmails = emails.map((e) =>
          e.id === email.id ? { ...e, read: false } : e
        );
        setEmails(revertedEmails);

        toast({
          title: "Error marking as read",
          description: response.error.message,
          variant: "destructive",
        });
      }
    }
    setSelectedEmail(email);
  };

  const handleArchive = (emailId: string) => {
    // For now, archiving is the same as deleting
    handleDelete(emailId);
  };

  const handleDelete = async (emailId: string) => {
    const originalEmails = [...emails]; // Store original state for potential revert
    const emailToDelete = emails.find((email) => email.id === emailId);

    // Optimistically update UI
    const updatedEmails = emails.filter((email) => email.id !== emailId);
    setEmails(updatedEmails);

    if (selectedEmail?.id === emailId) {
      setSelectedEmail(null); // Clear selection if deleted email was selected
    }

    const response = await deleteMessage(parseInt(emailId));

    if (response.error) {
      // Revert UI change if API call fails
      setEmails(originalEmails);
      // If the deleted email was selected, re-select it upon revert
      if (selectedEmail?.id === emailId) {
        setSelectedEmail(emailToDelete || null);
      }

      toast({
        title: "Error deleting message",
        description: response.error.message,
        variant: "destructive",
      });
      return;
    }

    // API call was successful
    toast({
      title: "Message deleted",
      description: "The message was successfully deleted.",
    });
  };

  const handleToggleRead = async (emailId: string) => {
    try {
      const response = await updateMessageReadStatus(parseInt(emailId));

      if (response.error) {
        toast({
          title: "Error updating read status",
          description: response.error.message,
          variant: "destructive",
        });
        return;
      }

      // Toggle the read status in the UI
      const updatedEmails = emails.map((email) =>
        email.id === emailId ? { ...email, read: !email.read } : email
      );
      setEmails(updatedEmails);

      // Update selected email if it's the one being toggled
      if (selectedEmail?.id === emailId) {
        setSelectedEmail({
          ...selectedEmail,
          read: !selectedEmail.read,
        });
      }

      toast({
        title: `Marked as ${
          updatedEmails.find((e) => e.id === emailId)?.read ? "read" : "unread"
        }`,
        description: "Email status updated successfully.",
      });
    } catch (error) {
      toast({
        title: "Error updating read status",
        description: "An unexpected error occurred.",
        variant: "destructive",
      });
    }
  };

  // Set current email (called from parent component)
  useEffect(() => {
    // This is where we would set up to listen for email address changes
    const handleStorageChange = () => {
      const storedEmail = localStorage.getItem(LOCAL_STORAGE_EMAIL_KEY);
      if (storedEmail && storedEmail !== currentEmail) {
        setCurrentEmail(storedEmail);
        fetchMessages(storedEmail);
        setSelectedEmail(null);
        setEmails([]);
      }
    };

    window.addEventListener("storage", handleStorageChange);
    return () => window.removeEventListener("storage", handleStorageChange);
  }, [currentEmail]);

  return (
    <Card className="h-full">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <div>
          <CardTitle className="text-xl">Inbox</CardTitle>
          <CardDescription>
            {currentEmail ? (
              <span className="flex items-center gap-2">
                {unreadCount > 0 ? (
                  <>
                    You have <Badge variant="secondary">{unreadCount}</Badge>{" "}
                    unread messages
                  </>
                ) : (
                  "No new messages"
                )}
              </span>
            ) : (
              "Generate an email to view messages"
            )}
          </CardDescription>
        </div>
        <Button
          variant="outline"
          size="icon"
          onClick={handleRefresh}
          disabled={refreshing || !currentEmail}
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
                      onToggleRead={() => handleToggleRead(email.id)}
                    />
                  ))
                ) : (
                  <div className="flex flex-col items-center justify-center py-12 text-center text-muted-foreground">
                    <h3 className="text-lg font-medium">No emails found</h3>
                    <p className="text-sm">
                      {searchQuery
                        ? "Try a different search term"
                        : currentEmail
                        ? "Your inbox is empty"
                        : "Generate an email address first"}
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
                        onToggleRead={() => handleToggleRead(email.id)}
                      />
                    ))
                ) : (
                  <div className="flex flex-col items-center justify-center py-12 text-center text-muted-foreground">
                    <h3 className="text-lg font-medium">No unread emails</h3>
                    <p className="text-sm">
                      {searchQuery
                        ? "Try a different search term"
                        : currentEmail
                        ? "You're all caught up!"
                        : "Generate an email address first"}
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
                  onToggleRead={() => handleToggleRead(selectedEmail.id)}
                />
              ) : (
                <div className="flex flex-col items-center justify-center h-full py-12 text-center text-muted-foreground">
                  <Clock className="h-12 w-12 mb-4" />
                  <h3 className="text-lg font-medium">No email selected</h3>
                  <p className="text-sm">
                    {currentEmail
                      ? "Select an email to view its contents"
                      : "Generate an email address first"}
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
