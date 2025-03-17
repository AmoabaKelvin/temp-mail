"use client";

import { Copy, RefreshCw, Shield } from "lucide-react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { toast } from "@/components/ui/use-toast";

import { generateRandomEmail } from "../lib/email-utils";

export function EmailHeader() {
  const [email, setEmail] = useState(generateRandomEmail());
  const [copying, setCopying] = useState(false);
  const [generating, setGenerating] = useState(false);

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
    setGenerating(true);
    setTimeout(() => {
      setEmail(generateRandomEmail());
      setGenerating(false);
      toast({
        title: "New email generated",
        description: "Your temporary email address has been updated.",
      });
    }, 500);
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
              value={email}
              readOnly
              className="w-full py-3 px-4 rounded-md focus:outline-none"
            />
            <Button
              variant="ghost"
              size="icon"
              className="absolute right-2"
              onClick={copyToClipboard}
              disabled={copying}
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
      </div>
    </div>
  );
}
