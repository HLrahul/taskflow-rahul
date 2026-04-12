import { useParams, useNavigate } from "react-router-dom";
import React, { useState } from "react";
import { ArrowLeft, ListChecks, Loader2 } from "lucide-react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import api from "@/lib/api";
import { CreateTaskDialog } from "@/components/CreateTaskDialog";
import { TaskCard, type Task } from "@/components/TaskCard";

interface ProjectData {
  id: string;
  name: string;
  description: string;
  owner_id: string;
  tasks: Task[];
}

interface StatsData {
  by_status: Record<string, number>;
  by_assignee: Record<string, number>;
}

export function ProjectDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user } = useAuth();

  const [draggedOverColumn, setDraggedOverColumn] = useState<string | null>(
    null,
  );
  const [filterStatus, setFilterStatus] = useState<string>("all");
  const [filterAssignee, setFilterAssignee] = useState<string>("all");

  const { data: team } = useQuery({
    queryKey: ["team"],
    queryFn: async () => {
      const res = await api.get("/team");
      return res.data.members;
    },
  });

  const {
    data: project,
    isLoading,
    error,
  } = useQuery<ProjectData>({
    queryKey: ["project", id],
    queryFn: async () => {
      const res = await api.get(`/projects/${id}`);
      return res.data;
    },
    enabled: !!id,
    refetchOnWindowFocus: false,
  });

  const { data: stats } = useQuery<StatsData>({
    queryKey: ["project-stats", id],
    queryFn: async () => {
      const res = await api.get(`/projects/${id}/stats`);
      return res.data;
    },
    enabled: !!id,
    refetchOnWindowFocus: false,
  });

  const updateMutation = useMutation({
    mutationFn: ({ taskId, payload }: { taskId: string; payload: any }) =>
      api.patch(`/tasks/${taskId}`, payload),
    // Optimistic update
    onMutate: async ({ taskId, payload }) => {
      await queryClient.cancelQueries({ queryKey: ["project", id] });

      const previousProject = queryClient.getQueryData<ProjectData>([
        "project",
        id,
      ]);

      queryClient.setQueryData<ProjectData>(["project", id], (old) => {
        if (!old) return old;
        return {
          ...old,
          tasks: old.tasks.map((t) =>
            t.id === taskId ? { ...t, status: payload.status } : t,
          ),
        };
      });

      return { previousProject };
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["project", id] });
      queryClient.invalidateQueries({ queryKey: ["project-stats", id] });
    },
    onError: (error: any, variables, context) => {
      // Revert optimistic update
      if (context?.previousProject) {
        queryClient.setQueryData(["project", id], context.previousProject);
      }

      if (error.response?.status === 403) {
        const taskTitle = variables.payload?.title || "this task";
        toast.error(`You are not authorized to move the task "${taskTitle}"`);
      } else {
        toast.error("Failed to move task safely");
      }
    },
  });

  const handleDragOver = (e: React.DragEvent, statusKey: string) => {
    e.preventDefault();
    if (draggedOverColumn !== statusKey) {
      setDraggedOverColumn(statusKey);
    }
  };

  const handleDragLeave = () => {
    setDraggedOverColumn(null);
  };

  const handleDrop = (e: React.DragEvent, statusKey: string) => {
    e.preventDefault();
    setDraggedOverColumn(null);

    const taskPayloadStr = e.dataTransfer.getData("taskPayload");
    if (!taskPayloadStr) return;

    const task = JSON.parse(taskPayloadStr);

    if (task.status !== statusKey) {
      updateMutation.mutate({
        taskId: task.id,
        payload: { ...task, status: statusKey },
      });
    }
  };

  if (isLoading)
    return (
      <div className="p-8 text-muted-foreground animate-pulse">
        Loading project details...
      </div>
    );
  if (error || !project)
    return (
      <div className="p-8 text-destructive">
        Failed to load project details.
      </div>
    );

  let filteredTasks = project.tasks || [];

  if (filterAssignee !== "all") {
    if (filterAssignee === "unassigned") {
      filteredTasks = filteredTasks.filter((t) => !t.assignee_id);
    } else {
      filteredTasks = filteredTasks.filter(
        (t) => t.assignee_id === filterAssignee,
      );
    }
  }

  const todoTasks = filteredTasks.filter((t) => t.status === "todo");
  const inProgressTasks = filteredTasks.filter(
    (t) => t.status === "in_progress",
  );
  const doneTasks = filteredTasks.filter((t) => t.status === "done");

  let columns = [
    {
      title: "To Do",
      key: "todo",
      tasks: todoTasks,
      color: "bg-slate-50 dark:bg-slate-900/50",
      border: "border-slate-200 dark:border-slate-800",
    },
    {
      title: "In Progress",
      key: "in_progress",
      tasks: inProgressTasks,
      color: "bg-blue-50 dark:bg-blue-950/20",
      border: "border-blue-100 dark:border-blue-900",
    },
    {
      title: "Done",
      key: "done",
      tasks: doneTasks,
      color: "bg-green-50 dark:bg-green-950/20",
      border: "border-green-100 dark:border-green-900",
    },
  ];

  if (filterStatus !== "all") {
    columns = columns.filter((c) => c.key === filterStatus);
  }

  const totalTasks =
    (stats?.by_status?.todo || 0) +
    (stats?.by_status?.in_progress || 0) +
    (stats?.by_status?.done || 0);

  return (
    <div className="flex-1 flex flex-col h-[calc(100vh-4rem)] overflow-hidden min-h-0 p-4 sm:p-6 lg:p-8 pt-4 sm:pt-6">
      {/* Header Area */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 shrink-0 mb-4 sm:mb-6">
        <div className="flex items-start gap-3 sm:gap-4 min-w-0">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => navigate("/")}
            className="cursor-pointer mt-1 sm:mt-3 hover:bg-accent hover:text-accent-foreground rounded-full w-8 h-8 sm:w-10 sm:h-10 shrink-0 -ml-1 sm:-ml-2"
          >
            <ArrowLeft className="h-4 w-4 sm:h-5 sm:w-5" />
          </Button>
          <div className="flex flex-col min-w-0">
            <h2 className="text-xl sm:text-2xl lg:text-3xl font-bold tracking-tight truncate">
              {project.name}
            </h2>
            {project.description && (
              <p className="text-xs sm:text-sm text-muted-foreground mt-1 sm:mt-2 max-w-2xl line-clamp-2">
                {project.description}
              </p>
            )}
          </div>
        </div>
        <CreateTaskDialog projectId={id!} />
      </div>

      {/* Stats Bar */}
      {stats && totalTasks > 0 && (
        <div className="flex flex-wrap gap-3 shrink-0 mb-4">
          <div className="flex items-center gap-2 rounded-lg border bg-background px-3 py-1.5 text-xs font-medium shadow-sm">
            <ListChecks className="h-3.5 w-3.5 text-muted-foreground" />
            <span className="text-muted-foreground">Total:</span>
            <span className="font-semibold">{totalTasks}</span>
          </div>
          <div className="flex items-center gap-2 rounded-lg border bg-background px-3 py-1.5 text-xs font-medium shadow-sm">
            <div className="h-2 w-2 rounded-full bg-slate-400" />
            <span className="text-muted-foreground">To Do:</span>
            <span className="font-semibold">{stats.by_status?.todo || 0}</span>
          </div>
          <div className="flex items-center gap-2 rounded-lg border bg-background px-3 py-1.5 text-xs font-medium shadow-sm">
            <Loader2 className="h-3 w-3 text-blue-500" />
            <span className="text-muted-foreground">In Progress:</span>
            <span className="font-semibold">
              {stats.by_status?.in_progress || 0}
            </span>
          </div>
          <div className="flex items-center gap-2 rounded-lg border bg-background px-3 py-1.5 text-xs font-medium shadow-sm">
            <div className="h-2 w-2 rounded-full bg-green-500" />
            <span className="text-muted-foreground">Done:</span>
            <span className="font-semibold">{stats.by_status?.done || 0}</span>
          </div>
        </div>
      )}

      {/* Filters Row */}
      <div className="flex flex-wrap gap-3 sm:gap-4 items-center mb-4 sm:mb-6 shrink-0">
        <Select value={filterStatus} onValueChange={setFilterStatus}>
          <SelectTrigger className="w-[140px] sm:w-[180px] bg-background text-xs sm:text-sm">
            <SelectValue placeholder="Filter by Status" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">All Statuses</SelectItem>
            <SelectItem value="todo">To Do</SelectItem>
            <SelectItem value="in_progress">In Progress</SelectItem>
            <SelectItem value="done">Done</SelectItem>
          </SelectContent>
        </Select>

        <Select value={filterAssignee} onValueChange={setFilterAssignee}>
          <SelectTrigger className="w-[140px] sm:w-[180px] bg-background text-xs sm:text-sm">
            <SelectValue placeholder="Filter by Assignee" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Any Assignee</SelectItem>
            <SelectItem value="unassigned">Unassigned</SelectItem>
            {user?.id && (
              <SelectItem value={user.id}>Assigned to Me</SelectItem>
            )}
            {team?.map((m: any) => (
              <SelectItem key={m.id} value={m.id}>
                {m.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {/* Kanban Board */}
      <div className="flex flex-1 gap-4 sm:gap-6 overflow-x-auto overflow-y-hidden pb-4 min-h-0">
        {columns.map((column) => (
          <div
            key={column.key}
            onDragOver={(e) => handleDragOver(e, column.key)}
            onDragLeave={handleDragLeave}
            onDrop={(e) => handleDrop(e, column.key)}
            className={`flex flex-col flex-1 min-w-[220px] sm:min-w-[250px] min-h-0 rounded-xl border transition-colors duration-200 ${draggedOverColumn === column.key ? "border-primary shadow-sm" : column.border} ${column.color}`}
          >
            {/* Column Header */}
            <div className="p-3 sm:p-4 flex items-center justify-between shrink-0 mb-2 border-b border-black/5 dark:border-white/5">
              <h3 className="font-semibold text-sm sm:text-base">
                {column.title}
              </h3>
              <div className="rounded-full bg-background/60 text-foreground px-2 sm:px-2.5 py-0.5 text-xs font-medium border border-black/5 dark:border-white/10 shadow-sm">
                {column.tasks.length}
              </div>
            </div>

            {/* Column Body / Tasks list */}
            <div className="flex-1 overflow-y-auto px-3 sm:px-4 pb-4 pt-2 sm:pt-4">
              {column.tasks.map((task) => (
                <TaskCard key={task.id} task={task} projectId={id!} />
              ))}

              <CreateTaskDialog
                projectId={id!}
                defaultStatus={column.key}
                trigger={
                  <div className="mt-1 flex flex-col gap-2 items-center justify-center text-xs sm:text-sm text-muted-foreground/60 border-2 border-dashed border-black/10 dark:border-white/10 hover:border-primary/50 hover:bg-accent/30 hover:text-foreground cursor-pointer transition-colors rounded-lg py-3">
                    {column.tasks.length > 0
                      ? "+ Add task"
                      : "Drop tasks here or click to add"}
                  </div>
                }
              />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
