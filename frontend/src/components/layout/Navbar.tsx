import { useNavigate } from "react-router-dom";
import { CheckCircle2 } from "lucide-react";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { ThemeToggle } from "@/components/ThemeToggle";
import { TeamManager } from "@/components/TeamManager";

export function Navbar() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <nav className="border-b sticky top-0 bg-background/95 backdrop-blur z-10 w-full px-3 sm:px-6 py-2.5 sm:py-3 flex items-center justify-between">
      <div
        className="flex items-center gap-2 font-bold text-lg sm:text-xl tracking-tight cursor-pointer"
        onClick={() => navigate("/")}
      >
        <div className="bg-primary text-primary-foreground p-1.5 rounded-lg flex items-center justify-center">
          <CheckCircle2 className="h-5 w-5" />
        </div>
        <span className="hidden xs:inline">TaskFlow</span>
      </div>

      <div className="flex items-center gap-2 sm:gap-6">
        <TeamManager />
        <ThemeToggle />
        <div className="flex items-center gap-2 sm:gap-3 pl-2 border-l">
          <Avatar className="h-7 w-7 sm:h-8 sm:w-8">
            <AvatarFallback className="bg-primary/10 text-primary font-medium text-xs sm:text-sm">
              {user?.name?.charAt(0).toUpperCase() || "U"}
            </AvatarFallback>
          </Avatar>
          <div className="text-sm font-medium hidden md:block">
            {user?.name}
          </div>
        </div>
        <Button
          variant="secondary"
          size="sm"
          onClick={handleLogout}
          className="cursor-pointer text-xs sm:text-sm"
        >
          Logout
        </Button>
      </div>
    </nav>
  );
}
