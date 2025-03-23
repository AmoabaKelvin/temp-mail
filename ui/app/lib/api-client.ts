// API client for interacting with the backend

// API base URL - can be configured per environment
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

// Types
interface ApiError {
  message: string;
  status: number;
}

interface ApiResponse<T> {
  data?: T;
  error?: ApiError;
}

export interface GeneratedAddress {
  email: string;
  expires_at: string;
}

export interface EmailMessage {
  id: number;
  from_address: string;
  subject: string;
  body: string;
  received_at: string;
}

// Generic fetch helper
async function apiFetch<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> {
  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
    });

    if (!response.ok) {
      return {
        error: {
          message: `API error: ${response.statusText}`,
          status: response.status,
        },
      };
    }

    const data = await response.json();
    return { data };
  } catch (error) {
    return {
      error: {
        message: error instanceof Error ? error.message : "Unknown error",
        status: 500,
      },
    };
  }
}

// API functions
export async function generateEmailAddress(): Promise<
  ApiResponse<GeneratedAddress>
> {
  return apiFetch<GeneratedAddress>("/v1/addresses", {
    method: "POST",
  });
}

export async function getMessages(
  email: string
): Promise<ApiResponse<EmailMessage[]>> {
  return apiFetch<EmailMessage[]>(
    `/v1/messages?email=${encodeURIComponent(email)}`
  );
}

export async function deleteMessage(
  id: number
): Promise<ApiResponse<{ message: string }>> {
  return apiFetch<{ message: string }>(`/v1/messages/${id}`, {
    method: "DELETE",
  });
}
