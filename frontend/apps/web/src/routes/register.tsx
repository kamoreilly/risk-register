import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import * as React from "react";

import { useAuth } from "@/hooks/useAuth";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/register")({
  component: RegisterComponent,
});

function RegisterComponent() {
  const navigate = useNavigate();
  const { register, isRegisterLoading, registerError, isAuthenticated } = useAuth();

  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [confirmPassword, setConfirmPassword] = React.useState("");
  const [name, setName] = React.useState("");
  const [validationError, setValidationError] = React.useState<string | null>(null);

  React.useEffect(() => {
    if (isAuthenticated) {
      navigate({ to: "/app" });
    }
  }, [isAuthenticated, navigate]);

  async function onSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setValidationError(null);

    if (password !== confirmPassword) {
      setValidationError("Passwords do not match");
      return;
    }

    if (password.length < 8) {
      setValidationError("Password must be at least 8 characters");
      return;
    }

    register({ email, password, name });
  }

  return (
    <div className="bg-background">
      <main className="grid min-h-svh place-items-center px-4">
        <Card className="w-full max-w-sm">
          <CardHeader className="border-b">
            <CardTitle>Create account</CardTitle>
            <CardDescription>Sign up to start managing risks.</CardDescription>
          </CardHeader>

          <CardContent>
            {(validationError || registerError) && (
              <div className="mb-4 rounded-md bg-destructive/10 p-3 text-sm text-destructive">
                {validationError || (registerError instanceof Error ? registerError.message : "Registration failed")}
              </div>
            )}
            <form className="grid gap-3" onSubmit={onSubmit}>
              <div className="grid gap-1.5">
                <Label htmlFor="name">Name</Label>
                <Input
                  id="name"
                  name="name"
                  type="text"
                  autoComplete="name"
                  placeholder="John Doe"
                  value={name}
                  onChange={(e) => setName(e.currentTarget.value)}
                  required
                />
              </div>

              <div className="grid gap-1.5">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  name="email"
                  type="email"
                  autoComplete="email"
                  placeholder="you@company.com"
                  value={email}
                  onChange={(e) => setEmail(e.currentTarget.value)}
                  required
                />
              </div>

              <div className="grid gap-1.5">
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  name="password"
                  type="password"
                  autoComplete="new-password"
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.currentTarget.value)}
                  required
                />
              </div>

              <div className="grid gap-1.5">
                <Label htmlFor="confirmPassword">Confirm password</Label>
                <Input
                  id="confirmPassword"
                  name="confirmPassword"
                  type="password"
                  autoComplete="new-password"
                  placeholder="••••••••"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.currentTarget.value)}
                  required
                />
              </div>

              <Button type="submit" disabled={isRegisterLoading}>
                {isRegisterLoading ? "Creating account..." : "Create account"}
              </Button>
            </form>
          </CardContent>

          <CardFooter className="justify-between">
            <Link to="/login" className={cn(buttonVariants({ variant: "link", size: "sm" }))}>
              Already have an account?
            </Link>
            <Link to="/" className={cn(buttonVariants({ variant: "outline", size: "sm" }))}>
              Back
            </Link>
          </CardFooter>
        </Card>
      </main>
    </div>
  );
}
