import { useNavigate, Link } from "react-router-dom";
import { CheckCircle2 } from "lucide-react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { useAuth } from "@/context/AuthContext";
import api from "@/lib/api";

const loginSchema = z.object({
  email: z.string().email("Please enter a valid email address."),
  password: z.string().min(1, "Password is required"),
});

export function Login() {
  const navigate = useNavigate();
  const { login } = useAuth();

  const form = useForm<z.infer<typeof loginSchema>>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: "", password: "" },
    reValidateMode: "onSubmit",
  });

  const onSubmit = async (values: z.infer<typeof loginSchema>) => {
    try {
      const response = await api.post("/auth/login", values);
      login(response.data.token, response.data.user);
      toast.success("Successfully logged in!");
      navigate("/");
    } catch (err: any) {
      toast.error(err.response?.data?.error || "Invalid login credentials.");
    }
  };

  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-4">
      <div className="flex items-center gap-2 font-bold text-2xl tracking-tight mb-8">
        <div className="bg-primary text-primary-foreground p-1.5 rounded-lg flex items-center justify-center">
          <CheckCircle2 className="h-6 w-6" />
        </div>
        TaskFlow
      </div>
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>Login to your account</CardTitle>
          <CardDescription>Enter your email below to login</CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-5" noValidate>
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input
                        type="email"
                        placeholder="m@example.com"
                        {...field}
                        onChange={(e) => {
                          form.clearErrors("email");
                          field.onChange(e);
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        {...field}
                        onChange={(e) => {
                          form.clearErrors("password");
                          field.onChange(e);
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button
                type="submit"
                className="w-full cursor-pointer"
                disabled={form.formState.isSubmitting}
              >
                {form.formState.isSubmitting ? "Logging in..." : "Login"}
              </Button>
            </form>
          </Form>
        </CardContent>
        <CardFooter className="flex flex-col gap-2 mt-2">
          <div className="text-sm text-center text-muted-foreground w-full">
            Don't have an account?{" "}
            <Link to="/register" className="text-primary hover:underline ml-1">
              Sign Up
            </Link>
          </div>
        </CardFooter>
      </Card>
    </div>
  );
}
