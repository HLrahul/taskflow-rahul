import { useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { MoreHorizontal, Pencil, Trash2 } from "lucide-react";
import { toast } from "sonner";
import React from "react";

import api from "@/lib/api";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { EditTaskDialog } from "./EditTaskDialog";

export interface Task {
  id: string;
  title: string;
  status: string;
  description: string;
  priority: string;
  project_id: string;
  assignee_id?: string;
  assignee_name?: string;
  assignee_email?: string;
}

export function TaskCard({
  task,
  projectId,
}: {
  task: Task;
  projectId: string;
}) {
  const queryClient = useQueryClient();
  const [editOpen, setEditOpen] = useState(false);

  const deleteMutation = useMutation({
    mutationFn: () => api.delete(`/tasks/${task.id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["project", projectId] });
      toast.success("Task deleted");
    },
    onError: (error: any) => {
      if (error.response?.status === 403) {
        toast.error("You don't have permission to delete this task");
      } else {
        toast.error("Failed to delete task");
      }
    },
  });

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case "high":
        return "destructive";
      case "medium":
        return "default";
      case "low":
        return "secondary";
      default:
        return "outline";
    }
  };

  const handleDragStart = (e: React.DragEvent) => {
    e.dataTransfer.setData("taskId", task.id);
    e.dataTransfer.setData("taskStatus", task.status);
    e.dataTransfer.setData("taskPayload", JSON.stringify(task));
    e.dataTransfer.effectAllowed = "move";

    // Add a slight transparency to the dragged card visually
    setTimeout(() => {
      const target = e.target as HTMLElement;
      if (target && target.style) {
        target.style.opacity = "0.4";
      }
    }, 0);
  };

  const handleDragEnd = (e: React.DragEvent) => {
    const target = e.target as HTMLElement;
    if (target && target.style) {
      target.style.opacity = "1";
    }
  };

  return (
    <>
      <Card
        draggable
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
        className="p-4 mb-3 shadow-sm hover:shadow-md transition-shadow bg-background cursor-grab active:cursor-grabbing border-border"
      >
        {/* Top row: content + menu */}
        <div className="flex items-start justify-between">
          <div className="flex flex-col flex-1 min-w-0 pr-2">
            <h4 className="text-sm font-semibold leading-tight text-foreground">
              {task.title}
            </h4>

            {task.description && (
              <p className="text-xs text-muted-foreground line-clamp-2 mt-1">
                {task.description}
              </p>
            )}
          </div>

          <div className="shrink-0 -mt-2 -mr-2">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="ghost"
                  size="icon"
                  className="h-7 w-7 mr-1 text-muted-foreground hover:bg-accent hover:text-accent-foreground"
                >
                  <MoreHorizontal className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-[150px]">
                <DropdownMenuItem
                  className="cursor-pointer text-xs"
                  onClick={() => setEditOpen(true)}
                >
                  <Pencil className="h-3.5 w-3.5 mr-2" />
                  Edit
                </DropdownMenuItem>
                <DropdownMenuItem
                  className="cursor-pointer text-xs text-destructive focus:text-destructive"
                  onClick={() => deleteMutation.mutate()}
                >
                  <Trash2 className="h-3.5 w-3.5 mr-2" />
                  Delete
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>

        {/* Bottom row: badge + avatar */}
        <div className="mt-1 flex items-center justify-between">
          <Badge
            variant={getPriorityColor(task.priority) as any}
            className="capitalize text-[10px] h-5 px-2.5"
          >
            {task.priority || "medium"}
          </Badge>

          {task.assignee_id && (
            <Tooltip>
              <TooltipTrigger asChild>
                <Avatar className="h-6 w-6 border cursor-default shadow-sm hover:ring-2 hover:ring-primary/40 transition-shadow">
                  <AvatarFallback className="text-[10px] bg-primary/10 text-primary uppercase font-medium">
                    {task.assignee_name ? task.assignee_name.charAt(0) : "U"}
                  </AvatarFallback>
                </Avatar>
              </TooltipTrigger>
              <TooltipContent
                align="end"
                className="flex flex-col gap-1 px-3 py-2 border shadow-md relative z-50"
              >
                <span className="font-semibold text-xs leading-none">
                  {task.assignee_name || "Unknown User"}
                </span>
                <span className="text-[10px] text-muted-foreground leading-none">
                  {task.assignee_email || "No email"}
                </span>
              </TooltipContent>
            </Tooltip>
          )}
        </div>
      </Card>

      <EditTaskDialog
        task={task}
        projectId={projectId}
        open={editOpen}
        onOpenChange={setEditOpen}
      />
    </>
  );
}
