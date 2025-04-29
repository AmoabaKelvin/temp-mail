"use client";

import { Copy, RefreshCw, Shield } from "lucide-react";
import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { toast } from "@/components/ui/use-toast";

import { generateEmailAddress } from "../lib/api-client";

export function EmailHeader() {
  const [email, setEmail] = useState<string>("");
  const [expiresAt, setExpiresAt] = useState<string>("");
  const [copying, setCopying] = useState(false);
  const [generating, setGenerating] = useState(false);
  const [loading, setLoading] = useState(true);

  // Fetch an email address on initial load
  useEffect(() => {
    // Check if we already have an email in localStorage
    const storedEmail = localStorage.getItem("currentEmail");

    if (storedEmail) {
      setEmail(storedEmail);
      setLoading(false);
      // We don't have expiry time, but that's okay for initial load
    } else {
      fetchEmailAddress();
    }
  }, []);

  const fetchEmailAddress = async () => {
    setGenerating(true);
    setLoading(true);
    const response = await generateEmailAddress();
    setGenerating(false);
    setLoading(false);

    if (response.error) {
      toast({
        title: "Error generating email",
        description: response.error.message,
        variant: "destructive",
      });
      return;
    }

    if (response.data) {
      setEmail(response.data.email);
      setExpiresAt(response.data.expires_at);

      // Save to localStorage for sharing with Inbox component
      localStorage.setItem("currentEmail", response.data.email);

      // Trigger a storage event for other components to detect
      window.dispatchEvent(new Event("storage"));
    }
  };

  const copyToClipboard = async () => {
    setCopying(true);
    try {
      await navigator.clipboard.writeText(email);
      toast({
        title: "Email copied",
        description: "The email address has been copied to your clipboard.",
      });
    } catch (err) {
      toast({
        title: "Failed to copy",
        description: "Could not copy the email address to clipboard.",
        variant: "destructive",
      });
    } finally {
      setCopying(false);
    }
  };

  const generateNewEmail = () => {
    fetchEmailAddress();
    toast({
      title: "Generating new email",
      description: "Your temporary email address is being updated.",
    });
  };

  // Format expiration time
  const formatExpiryTime = () => {
    if (!expiresAt) return "";
    const expiresAtDate = new Date(expiresAt);
    return expiresAtDate.toLocaleString();
  };

  return (
    <div className="flex flex-col items-center text-center">
      <div className="inline-flex items-center gap-2 mb-2 text-primary">
        <Shield className="h-5 w-5" />
        <span className="text-sm font-medium">
          Open Source Disposable Email
        </span>
      </div>
      <h1 className="text-5xl font-bold text-gray-900 mb-4">
        Free Temporary Email
      </h1>
      <p className="text-lg text-gray-600 max-w-3xl mb-8">
        Receive emails anonymously with our free temporary email address
        generator.
      </p>

      <div className="w-full max-w-3xl bg-gray-50 border border-gray-100 rounded-lg p-6">
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="relative flex-1 flex items-center border rounded-md bg-white">
            <input
              type="text"
              value={loading ? "Loading..." : email}
              readOnly
              className="w-full py-3 px-4 rounded-md focus:outline-none"
            />
            <Button
              variant="ghost"
              size="icon"
              className="absolute right-2"
              onClick={copyToClipboard}
              disabled={copying || loading}
            >
              <Copy className="h-5 w-5 text-gray-500" />
              <span className="sr-only">Copy email address</span>
            </Button>
          </div>

          <Button
            variant="outline"
            className="py-6 px-6 flex items-center justify-center gap-2 bg-white"
            onClick={generateNewEmail}
            disabled={generating}
          >
            <RefreshCw
              className={`h-4 w-4 ${generating ? "animate-spin" : ""}`}
            />
            Change email
          </Button>
        </div>
        {expiresAt && (
          <div className="mt-2 text-sm text-gray-500">
            Expires at: {formatExpiryTime()}
          </div>
        )}
      </div>
    </div>
  );
}
