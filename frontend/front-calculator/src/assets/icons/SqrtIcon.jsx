export function SqrtIcon({ size = 24, className }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      role="presentation"
      aria-hidden="true"
    >
      <path d="M4 12h3" />
      <path d="M7 12l2.5 6L15 5" />
      <path d="M13 5h7" />
    </svg>
  )
}
