import {Card, CardContent, CardHeader, CardTitle} from './ui/card';
import * as React from 'react';
import {Skeleton} from './ui/skeleton';

export function EventOverviewSkeleton() {
  const overviewTotalSkeleton = []
  const chartTotalSkeleton = []

  for (let i = 0; i < 4; i++) {
    overviewTotalSkeleton.push(
      <Card key={i}>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <Skeleton className="h-6 w-[100px]" />
          <Skeleton className="h-6 w-[25px]" />
        </CardHeader>
        <CardContent>
          <Skeleton className="mt-2 h-8 w-[50px]" />
          <Skeleton className="mt-4 h-2 w-full" />
          <Skeleton className="mt-1 h-2 w-full" />
          <Skeleton className="mt-1 h-2 w-[80px]" />
        </CardContent>
      </Card>
    )
  }

  for (let i = 0; i < 7; i++) {
    chartTotalSkeleton.push(
      <div key={i} className="flex flex-col items-center">
        <Skeleton className="h-[65%] w-16 mb-6" />
        <Skeleton className="h-8 w-12 mb-8" />
      </div>
    )
  }

  return (<>
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {overviewTotalSkeleton}
    </div>
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
      <Card className="col-span-4 p-4 h-full">
        <CardHeader>
          <CardTitle>
            <Skeleton className="h-6 w-[50%]" />
          </CardTitle>
        </CardHeader>
        <CardContent className="h-full flex flex-row space-x-6">
          <div className="flex flex-col">
            <Skeleton className="h-8 w-8 mb-8" />
            <Skeleton className="h-8 w-8 mb-8" />
            <Skeleton className="h-8 w-8 mb-8" />
            <Skeleton className="h-8 w-8 mb-8" />
          </div>
          {chartTotalSkeleton}
        </CardContent>
      </Card>
      <Card className="col-span-3 p-4">
        <CardHeader>
          <CardTitle>
            <Skeleton className="h-6 w-[150px] mb-2" />
            <Skeleton className="h-4 w-[50%]" />
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-full" />
          </div>
        </CardContent>
      </Card>
    </div>
  </>)
}