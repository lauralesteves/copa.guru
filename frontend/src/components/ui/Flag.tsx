interface FlagProps {
  code: string;
  size?: 48 | 96;
  className?: string;
}

export function Flag({ code, size = 48, className = '' }: FlagProps) {
  return (
    <img
      src={`/images/flags/${size}/${code}.png`}
      alt={code}
      width={size === 48 ? 40 : 64}
      height={size === 48 ? 30 : 48}
      className={`inline-block flag-wave ${className}`}
      loading="lazy"
    />
  );
}
