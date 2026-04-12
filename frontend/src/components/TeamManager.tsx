import { useState } from "react";
import { Users, Loader2 } from "lucide-react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import api from "@/lib/api";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";

export function TeamManager() {
  const [open, setOpen] = useState(false);
  const [email, setEmail] = useState("");
  const queryClient = useQueryClient();

  const { data: team, isLoading } = useQuery({
    queryKey: ["team"],
    queryFn: async () => {
      const res = await api.get("/team");
      return res.data.members;
    },
    enabled: open,
  });

  const mutation = useMutation({
    mutationFn: async (userEmail: string) => {
      const res = await api.post("/team", { email: userEmail });
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["team"] });
      toast.success("Team member added!");
      setEmail("");
    },
    onError: (err: any) => {
      toast.error(err.response?.data?.error || "Failed to add user");
    },
  });

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="ghost" size="sm" className="gap-2 shrink-0 cursor-pointer text-muted-foreground hover:bg-accent hover:text-foreground">
          <Users className="h-4 w-4" />
          <span className="hidden sm:inline font-medium">Team</span>
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Manage Team</DialogTitle>
        </DialogHeader>
        <div className="flex flex-col gap-4 py-2">
          <div className="flex gap-2">
            <Input
              placeholder="user@example.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter" && email) {
                  e.preventDefault();
                  mutation.mutate(email);
                }
              }}
            />
            <Button 
              onClick={() => mutation.mutate(email)} 
              disabled={!email || mutation.isPending}
              className="cursor-pointer shrink-0"
            >
              {mutation.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : "Add"}
            </Button>
          </div>

          <div className="mt-4 space-y-3">
            <h4 className="text-sm font-medium text-muted-foreground">Current Members</h4>
            {isLoading ? (
              <div className="text-sm text-muted-foreground">Loading members...</div>
            ) : team?.length > 0 ? (
              <div className="max-h-[300px] overflow-y-auto space-y-2 pr-2">
                {team.map((member: any) => (
                  <div key={member.id} className="flex items-center gap-3 p-2 rounded-md border bg-accent/20">
                    <Avatar className="h-8 w-8">
                      <AvatarFallback className="bg-primary/10 text-primary text-xs">
                        {member.name.charAt(0).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                    <div className="flex flex-col">
                      <span className="text-sm font-medium">{member.name}</span>
                      <span className="text-xs text-muted-foreground">{member.email}</span>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-sm border-dashed border-2 p-6 text-center text-muted-foreground rounded-md">
                No team members yet. Invite someone!
              </div>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
