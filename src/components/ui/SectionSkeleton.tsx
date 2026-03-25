interface SectionSkeletonProps {
  lines?: number;
}

export function SectionSkeleton({ lines = 3 }: SectionSkeletonProps) {
  return (
    <div className="py-20 px-6">
      <div className="max-w-7xl mx-auto">
        <div className="skeleton h-10 w-64 mx-auto mb-4" />
        <div className="skeleton h-4 w-48 mx-auto mb-12" />
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {Array.from({ length: lines }, (_, i) => (
            <div
              key={`skel-${i}`}
              className="skeleton h-64 rounded-xl"
            />
          ))}
        </div>
      </div>
    </div>
  );
}
