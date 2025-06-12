import React from 'react';

interface TypographyProps {
  children: React.ReactNode;
  variant?: 'display-large' | 'display-medium' | 'heading-xsmall' | 'text-large-paragraph';
  color?: 'dark' | 'light' | 'light-gray-4' | 'white' | 'link-blue';
  className?: string;
}

export const Typography: React.FC<TypographyProps> = ({ 
  children, 
  variant = 'text-large-paragraph', 
  color = 'dark',
  className = '' 
}) => {
  const colorClasses = {
    dark: 'text-[#333333]',
    light: 'text-[#757575]', 
    'light-gray-4': 'text-[#e0e0e0]',
    white: 'text-[#ffffff]',
    'link-blue': 'text-[#2d7ff9]'
  };

  const variantClasses = {
    'display-large': 'display-large',
    'display-medium': 'display-medium', 
    'heading-xsmall': 'heading-xsmall',
    'text-large-paragraph': 'text-large-paragraph'
  };

  return (
    <div className={`${variantClasses[variant]} ${colorClasses[color]} ${className}`}>
      {children}
    </div>
  );
};

interface FontDemoProps {
  className?: string;
}

export const FontDemo: React.FC<FontDemoProps> = ({ className = '' }) => {
  return (
    <div className={`bg-white p-16 ${className}`}>
      <div className="flex flex-col items-center w-full">
        <div className="flex flex-col gap-10 items-center justify-start p-16 w-full">
          {/* Header */}
          <div className="w-full">
            <div className="flex flex-col gap-6 items-start justify-start w-full">
              <Typography variant="display-medium" color="dark">
                Fonts
              </Typography>
              <div className="bg-[#e0e0e0] h-0.5 w-full" />
            </div>
          </div>
          
          {/* Font Scale and Description */}
          <div className="w-full">
            <div className="flex flex-row items-center justify-between w-full">
              {/* Font Scale */}
              <div className="flex flex-col gap-8 items-start">
                <Typography variant="display-large" color="dark" className="font-bold tracking-[-0.312px]">
                  SF Pro Display
                </Typography>
                <Typography variant="display-large" color="dark" className="font-normal tracking-[-2.184px]">
                  SF Pro Text
                </Typography>
              </div>
              
              {/* Description */}
              <div className="w-[280px]">
                <div className="flex flex-col gap-2 items-center justify-start w-[280px]">
                  <Typography variant="heading-xsmall" color="dark" className="w-[280px]">
                    Download
                  </Typography>
                  <div className="w-[280px]">
                    <Typography variant="text-large-paragraph" color="light" className="mb-4">
                      Make sure you have the right version of San Francisco Pro Display & Text installed.
                    </Typography>
                    <Typography variant="text-large-paragraph" color="light">
                      You can download San Francisco from:{' '}
                      <a 
                        href="https://developer.apple.com/fonts/" 
                        className="text-[#2d7ff9] underline cursor-pointer hover:text-[#1283da]"
                      >
                        https://developer.apple.com/fonts/
                      </a>
                    </Typography>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}; 