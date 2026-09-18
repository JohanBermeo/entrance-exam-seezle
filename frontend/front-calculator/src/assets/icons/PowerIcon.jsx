export function PowerIcon({ size = 24, className }) {
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
      <path d="M4 15l3-3 3 3" />
      <path d="M7 12v6" />
      <path d="M13 13l2.5-2.5L18 13" />
      <path d="M15.5 10.5V18" />
    </svg>
  )
}
