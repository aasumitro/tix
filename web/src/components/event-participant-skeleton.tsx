import {Skeleton} from './ui/skeleton';
import * as React from 'react';

export function EventParticipantSkeleton() {
  const skeleton = []

  for (let i = 0; i < 5; i++) {
    skeleton.push(<div key={i} className="flex flex-row items-center bg-gray-100 p-4 rounded-md">
      <Skeleton className="h-8 w-10 mr-4 bg-gray-200" />
      <Skeleton className="h-8 w-full bg-gray-200" />
    </div>)
  }

  return (
    <>
      <div className="flex flex-row items-center justify-between mb-10">
        <div>
          <Skeleton className="w-36 h-8"/>
          <Skeleton className="w-[400px] h-4 mt-2"/>
        </div>
        <Skeleton className="w-52 h-12 mt-2"/>
      </div>

      <div className="flex items-center justify-between gap-2">
        <div className="flex flex-1 items-center space-x-2">
          <Skeleton className="h-10 w-[400px]" />
        </div>
        <Skeleton className="h-10 w-[110px] mr-3" />
        <Skeleton className="h-10 w-[50px]" />
      </div>

      <div className="w-full h-[450px] bg-gray-50 rounded-xl flex flex-col gap-4 p-8 animate-pulse">
        {skeleton}
      </div>

      <div className="flex items-center justify-between gap-2">
        <div className="flex flex-1 items-center space-x-2">
          <Skeleton className="h-4 w-[150px]" />
        </div>
        <Skeleton className="h-8 w-[80px] mr-3" />
        <Skeleton className="h-8 w-[50px]" />
      </div>
    </>
  )
}