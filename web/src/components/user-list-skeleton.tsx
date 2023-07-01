import * as React from 'react';
import {Skeleton} from './ui/skeleton';
import {Card, CardContent} from './ui/card';

export function UserListSkeleton() {
  const skeleton = [];

  for (let i = 0; i < 5; i++) {
    skeleton.push(
      <Card key={i} className="rounded-md">
        <CardContent className="flex flex-row justify-between items-center px-4 py-3">
          <div className="flex">
            <Skeleton className="h-8 w-8 rounded-full" />
            <div className="flex flex-row items-center gap-2 ml-4">
              <Skeleton className="h-4 w-20" />
              <Skeleton className="h-3 w-40" />
            </div>
          </div>
          <div className="flex gap-2">
            <Skeleton className="h-8 w-20 rounded-xl" />
            <Skeleton className="h-8 w-20 rounded-xl" />
          </div>
          <div className="flex items-center gap-2">
            <Skeleton className="h-6 w-6" />
            <Skeleton className="h-6 w-6" />
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="container w-full h-full space-y-4 py-12 flex flex-col">
      <div className="flex flex-row items-center justify-between mb-10">
        <div>
          <Skeleton className="w-20 h-8"/>
          <Skeleton className="w-80 h-4 mt-2"/>
        </div>
        <Skeleton className="w-40 h-12 mt-2"/>
      </div>

      <div className="flex items-center justify-between gap-2">
        <div className="flex flex-1 items-center space-x-2">
         <Skeleton className="h-8 w-[250px]" />
        </div>
        <Skeleton className="h-8 w-[40px]" />
        <Skeleton className="h-8 w-[40px]" />
      </div>

      <div className="flex flex-col pt-4 gap-4">
        {skeleton}
      </div>
    </div>
  )
}