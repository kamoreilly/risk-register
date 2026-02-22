import { render, screen, fireEvent } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { LoginComponent } from "@/routes/login";

// Mocks
const mockNavigate = vi.fn();
const mockLogin = vi.fn();

vi.mock("@tanstack/react-router", () => ({
  useNavigate: () => mockNavigate,
  Link: ({ to, children }: { to: string; children: React.ReactNode }) => <a href={to}>{children}</a>,
}));

vi.mock("@/hooks/useAuth", () => ({
  useAuth: () => ({
    login: mockLogin,
    isLoginLoading: false,
    loginError: null,
    isAuthenticated: false,
  }),
}));

describe("LoginComponent", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders login form", () => {
    render(<LoginComponent />);
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /sign in/i })).toBeInTheDocument();
  });

  it("submits form with credentials", async () => {
    render(<LoginComponent />);
    
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: "test@example.com" } });
    fireEvent.change(screen.getByLabelText(/password/i), { target: { value: "password" } });
    
    fireEvent.click(screen.getByRole("button", { name: /sign in/i }));
    
    expect(mockLogin).toHaveBeenCalledWith({ email: "test@example.com", password: "password" });
  });
});
