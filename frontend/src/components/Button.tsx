import React from 'react';

export interface ButtonProps {
  children: React.ReactNode;
  variant?: 'default' | 'primary' | 'danger' | 'secondary' | 'text';
  size?: 'small' | 'default' | 'large' | 'xlarge';
  icon?: React.ReactNode;
  iconOnly?: boolean;
  className?: string;
  onClick?: () => void;
  disabled?: boolean;
}

export const Button: React.FC<ButtonProps> = ({
  children,
  variant = 'default',
  size = 'default',
  icon,
  iconOnly = false,
  className = '',
  onClick,
  disabled = false,
}) => {
  const baseClasses = 'relative inline-flex items-center justify-center font-medium transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2';
  
  // Size classes for regular buttons
  const sizeClasses = {
    small: 'px-2.5 py-1.5 text-[11px] leading-4 rounded-[3px] gap-1.5',
    default: 'px-3 py-2 text-[13px] leading-4 rounded-[3px] gap-2',
    large: 'px-4 py-2.5 text-[15px] leading-5 rounded-[3px] gap-2.5',
    xlarge: 'px-5 py-3 text-[17px] leading-6 rounded-[3px] gap-3',
  };

  // Text button size classes
  const textSizeClasses = {
    small: 'text-[11px] leading-4',
    default: 'text-[13px] leading-4', 
    large: 'text-[15px] leading-5',
    xlarge: 'text-[17px] leading-6',
  };

  // Icon only size classes
  const iconOnlySizeClasses = {
    small: 'p-1.5 rounded-[3px]',
    default: 'p-2 rounded-[3px]',
    large: 'p-2.5 rounded-[3px]',
    xlarge: 'p-3 rounded-[3px]',
  };

  // Variant classes for regular buttons
  const variantClasses = {
    default: 'bg-[#f2f2f2] text-[#333333] border border-[#e0e0e0] hover:bg-[#e0e0e0] focus:ring-[#757575]',
    primary: 'bg-[#2d7ff9] text-white border border-[#2d7ff9] hover:bg-[#1e6ee8] focus:ring-[#2d7ff9]',
    danger: 'bg-[#ef3061] text-white border border-[#ef3061] hover:bg-[#d9285a] focus:ring-[#ef3061]',
    secondary: 'bg-white text-[#333333] border border-[#e0e0e0] hover:bg-[#fafafa] focus:ring-[#757575]',
    text: 'bg-transparent text-[#2d7ff9] hover:text-[#1e6ee8] focus:ring-[#2d7ff9]',
  };

  // Icon size classes
  const iconSizeClasses = {
    small: 'w-3 h-3',
    default: 'w-3.5 h-3.5',
    large: 'w-4 h-4',
    xlarge: 'w-5 h-5',
  };

  // Disabled classes
  const disabledClasses = 'opacity-50 cursor-not-allowed';

  const getSizeClass = () => {
    if (iconOnly) return iconOnlySizeClasses[size];
    if (variant === 'text') return textSizeClasses[size];
    return sizeClasses[size];
  };

  const combinedClasses = `
    ${baseClasses}
    ${getSizeClass()}
    ${variantClasses[variant]}
    ${disabled ? disabledClasses : ''}
    ${className}
  `.trim().replace(/\s+/g, ' ');

  return (
    <button
      className={combinedClasses}
      onClick={onClick}
      disabled={disabled}
    >
      {icon && (
        <span className={iconSizeClasses[size]}>
          {icon}
        </span>
      )}
      {!iconOnly && children}
    </button>
  );
};

// Icon component for edit icon
export const EditIcon: React.FC<{ className?: string }> = ({ className = '' }) => (
  <svg
    className={`block ${className}`}
    fill="none"
    viewBox="0 0 12 12"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      clipRule="evenodd"
      d="M0.31216 8.07501L5.81284 2.57433C5.98556 2.40161 6.26443 2.40156 6.43643 2.5735L8.43857 4.57497C8.60887 4.74521 8.61024 5.02607 8.43784 5.1985L2.93716 10.7C2.76444 10.8727 2.43093 11.0122 2.19153 11.0122H0.433465C0.188784 11.0122 0 10.8181 0 10.5787V8.82064C0 8.57595 0.139759 8.24741 0.31216 8.07501ZM9.37212 0.259609L10.7534 1.64043C11.0947 1.9816 11.0971 2.54082 10.7535 2.88412L9.9298 3.70697C9.76221 3.87439 9.48539 3.87226 9.31339 3.70032L7.31125 1.69885C7.14096 1.5286 7.1363 1.25104 7.30465 1.08269L8.12791 0.259432C8.47639 -0.0890423 9.02844 -0.0839613 9.37212 0.259609Z"
      fill="currentColor"
      fillRule="evenodd"
    />
  </svg>
);

// Tooltip component for icon-only buttons
export interface TooltipProps {
  children: React.ReactNode;
  text: string;
  className?: string;
}

export const Tooltip: React.FC<TooltipProps> = ({ children, text, className = '' }) => {
  return (
    <div className={`relative inline-block group ${className}`}>
      {children}
      <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none">
        <div className="bg-[#333333] text-white text-[13px] px-1.5 py-2 rounded-[3px] whitespace-nowrap">
          {text}
        </div>
      </div>
    </div>
  );
}; 