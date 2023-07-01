import {cn} from '../../libs/utils';
import React from 'react';
import {useLocation, useNavigate} from 'react-router-dom';
export function MainNav({
  className,
  ...props
}: React.HTMLAttributes<HTMLElement>) {
  const navigate = useNavigate()
  const eventID = localStorage.getItem("current_event")
  const location = useLocation()

  return (
    <nav
      className={cn("flex items-center space-x-4 lg:space-x-6", className)}
      {...props}
    >
      <button
        onClick={() => navigate(`/event/overview/${eventID}`)}
        className={`
          text-sm font-medium transition-colors hover:text-primary 
          ${location.pathname.includes("overview") ? "text-black" : "text-muted-foreground "}
        `}
      >
        Overview
      </button>
      <button
        onClick={() => navigate(`event/participants/${eventID}`)}
        className={`
          text-sm font-medium transition-colors hover:text-primary 
          ${location.pathname.includes("participants") ? "text-black" : "text-muted-foreground "}
        `}
      >
        Participants
      </button>
      <button
        onClick={() => navigate('/users')}
        className={`
          text-sm font-medium transition-colors hover:text-primary 
          ${location.pathname.includes("users") ? "text-black" : "text-muted-foreground "}
        `}
      >
        Users
      </button>
    </nav>
  )
}
