import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import {
  Folders,
  CheckCircle2,
  Clock,
  ListTodo,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";

import api from "@/lib/api";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { CreateProjectDialog } from "@/components/CreateProjectDialog";

interface Project {
  id: string;
  name: string;
  description: string;
  created_at: string;
}

interface ProjectsResponse {
  projects: Project[];
  total: number;
  page: number;
  limit: number;
}

interface StatsData {
  by_status: Record<string, number>;
  by_assignee: Record<string, number>;
}

function ProjectStatsBar({ projectId }: { projectId: string }) {
  const { data: stats } = useQuery<StatsData>({
    queryKey: ["project-stats", projectId],
    queryFn: async () => {
      const res = await api.get(`/projects/${projectId}/stats`);
      return res.data;
    },
    staleTime: 30000,
  });

  if (!stats) return null;

  const todo = stats.by_status?.todo || 0;
  const inProgress = stats.by_status?.in_progress || 0;
  const done = stats.by_status?.done || 0;
  const total = todo + inProgress + done;

  if (total === 0) return null;

  const donePercent = Math.round((done / total) * 100);

  return (
    <div className="space-y-2.5">
      {/* Progress bar */}
      <div className="w-full h-1.5 bg-muted rounded-full overflow-hidden">
        <div
          className="h-full bg-green-500 rounded-full transition-all duration-500"
          style={{ width: `${donePercent}%` }}
        />
      </div>

      {/* Counts */}
      <div className="flex items-center gap-3 text-[11px] text-muted-foreground">
        <span className="flex items-center gap-1">
          <ListTodo className="h-3 w-3" />
          {todo}
        </span>
        <span className="flex items-center gap-1">
          <Clock className="h-3 w-3 text-blue-500" />
          {inProgress}
        </span>
        <span className="flex items-center gap-1">
          <CheckCircle2 className="h-3 w-3 text-green-500" />
          {done}
        </span>
        <span className="ml-auto font-medium text-foreground/70">
          {donePercent}%
        </span>
      </div>
    </div>
  );
}

export function Dashboard() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);

  const { data, isLoading, error } = useQuery<ProjectsResponse>({
    queryKey: ["projects", page, limit],
    queryFn: async () => {
      const res = await api.get(`/projects?page=${page}&limit=${limit}`);
      return res.data;
    },
    refetchOnWindowFocus: false,
  });

  const projects = data?.projects || [];
  const total = data?.total || 0;
  const totalPages = Math.ceil(total / limit);

  const handleLimitChange = (val: string) => {
    setLimit(Number(val));
    setPage(1); // reset to first page on limit change
  };

  return (
    <div className="flex-1 flex flex-col p-4 sm:p-6 lg:p-8 pt-4 sm:pt-6">
      <div className="flex items-center justify-between gap-4 mb-6">
        <h2 className="text-2xl sm:text-3xl font-bold tracking-tight">
          Dashboard
        </h2>
        <div className="flex items-center gap-3">
          {/* Limit selector */}
          <div className="flex items-center gap-2">
            <span className="text-xs text-muted-foreground hidden sm:inline">
              Show:
            </span>
            <Select value={String(limit)} onValueChange={handleLimitChange}>
              <SelectTrigger className="w-[70px] h-9 text-xs">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="5">5</SelectItem>
                <SelectItem value="10">10</SelectItem>
                <SelectItem value="15">15</SelectItem>
                <SelectItem value="20">20</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <CreateProjectDialog />
        </div>
      </div>

      {isLoading ? (
        <div className="text-muted-foreground mt-10 text-sm">
          Loading projects...
        </div>
      ) : error ? (
        <div className="text-destructive mt-10">Failed to load projects.</div>
      ) : projects.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-md border border-dashed p-12 sm:p-24 text-center mt-10">
          <Folders className="h-10 w-10 text-muted-foreground mb-4 opacity-20" />
          <h3 className="text-lg sm:text-xl font-bold tracking-tight text-foreground/80">
            No projects yet
          </h3>
          <p className="text-sm text-muted-foreground mb-6">
            Create your first project to start managing tasks.
          </p>
          <CreateProjectDialog />
        </div>
      ) : (
        <>
          <div className="grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 flex-1">
            {projects.map((project) => (
              <Card
                key={project.id}
                className="cursor-pointer hover:border-primary/50 hover:bg-accent/10 transition-all flex flex-col"
                onClick={() => navigate(`/project/${project.id}`)}
              >
                <CardHeader className="pb-0 flex-1">
                  <div className="w-[95%]">
                    <CardTitle className="text-base sm:text-lg font-semibold leading-tight">
                      {project.name}
                    </CardTitle>
                    <p className="text-xs sm:text-sm text-muted-foreground line-clamp-3 mt-1.5">
                      {project.description || "No description provided."}
                    </p>
                  </div>
                </CardHeader>
                <CardContent className="pt-4">
                  <ProjectStatsBar projectId={project.id} />
                </CardContent>
              </Card>
            ))}
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-center gap-1 mt-8 pt-4 border-t">
              <Button
                variant="outline"
                size="icon"
                className="h-8 w-8"
                disabled={page <= 1}
                onClick={() => setPage((p) => Math.max(1, p - 1))}
              >
                <ChevronLeft className="h-4 w-4" />
              </Button>

              {Array.from({ length: totalPages }, (_, i) => i + 1).map(
                (pageNum) => (
                  <Button
                    key={pageNum}
                    variant={pageNum === page ? "default" : "outline"}
                    size="sm"
                    className="h-8 w-8 text-xs p-0"
                    onClick={() => setPage(pageNum)}
                  >
                    {pageNum}
                  </Button>
                ),
              )}

              <Button
                variant="outline"
                size="icon"
                className="h-8 w-8"
                disabled={page >= totalPages}
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              >
                <ChevronRight className="h-4 w-4" />
              </Button>

              <span className="text-xs text-muted-foreground ml-3">
                {total} project{total !== 1 ? "s" : ""}
              </span>
            </div>
          )}
        </>
      )}
    </div>
  );
}
