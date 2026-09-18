export function PlusMinusIcon({ size = 24, className }) {
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
      <path d="M4 8h16" />
      <path d="M12 12v8" />
      <path d="M8 16h8" />
    </svg>
  )
}
