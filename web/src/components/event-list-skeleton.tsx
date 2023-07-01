import {Card, CardContent, CardFooter, CardHeader, CardTitle} from './ui/card';
import {Skeleton} from './ui/skeleton';
import * as React from 'react';

interface EventSkeletonProps {
  display: number;
}

export function EventListSkeleton(props: EventSkeletonProps) {
  const skeleton = [];

  for (let i = 0; i < props.display; i++) {
    skeleton.push(
        <Card key={i} className="w-60 h-60 rounded-lg text-left">
          <CardHeader className="space-y-0 pb-2">
            <CardTitle className="text-md font-medium">
              <Skeleton className="h-4 w-[75px]" />
            </CardTitle>
          </CardHeader>
          <CardContent className="mt-8">
            <div className="flex gap-2 items-center">
              <Skeleton className="h-8 w-8" />
              <Skeleton className="h-8 w-12" />
            </div>
            <Skeleton className="h-4 w-24 mt-2" />
          </CardContent>
          <CardFooter className="text-sm font-normal flex items-baseline flex-row gap-2 mt-4">
            <Skeleton className="h-4 w-4" />
            <Skeleton className="h-4 w-32" />
          </CardFooter>
        </Card>
    )
  }

  return (<div className="flex flex-wrap gap-4 py-24 justify-center"> {skeleton} </div>)
}