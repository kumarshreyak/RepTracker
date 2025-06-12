import React from 'react';

interface TypographyProps {
  children: React.ReactNode;
  variant?: 
    | 'text-small' 
    | 'text-default' 
    | 'text-large' 
    | 'text-xlarge'
    | 'text-small-paragraph'
    | 'text-default-paragraph'
    | 'text-large-paragraph'
    | 'text-xlarge-paragraph'
    | 'heading-xsmall'
    | 'heading-small'
    | 'heading-default'
    | 'heading-large'
    | 'heading-xlarge'
    | 'heading-xxlarge'
    | 'heading-xsmall-caps'
    | 'heading-small-caps'
    | 'heading-default-caps';
  color?: 'dark' | 'light' | 'light-gray-4' | 'white' | 'link-blue' | 'red';
  className?: string;
}

export const Typography: React.FC<TypographyProps> = ({ 
  children, 
  variant = 'text-default', 
  color = 'dark',
  className = '' 
}) => {
  const colorClasses = {
    dark: 'text-[#333333]',
    light: 'text-[#757575]', 
    'light-gray-4': 'text-[#e0e0e0]',
    white: 'text-[#ffffff]',
    'link-blue': 'text-[#2d7ff9]',
    red: 'text-[#ef3061]'
  };

  const variantClasses = {
    'text-small': 'text-small',
    'text-default': 'text-default',
    'text-large': 'text-large',
    'text-xlarge': 'text-xlarge',
    'text-small-paragraph': 'text-small-paragraph',
    'text-default-paragraph': 'text-default-paragraph',
    'text-large-paragraph': 'text-large-paragraph',
    'text-xlarge-paragraph': 'text-xlarge-paragraph',
    'heading-xsmall': 'heading-xsmall',
    'heading-small': 'heading-small',
    'heading-default': 'heading-default',
    'heading-large': 'heading-large',
    'heading-xlarge': 'heading-xlarge',
    'heading-xxlarge': 'heading-xxlarge',
    'heading-xsmall-caps': 'heading-xsmall-caps',
    'heading-small-caps': 'heading-small-caps',
    'heading-default-caps': 'heading-default-caps'
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
        <div className="flex flex-col gap-16 items-center justify-start p-16 w-full">
          {/* Header */}
          <div className="w-full">
            <div className="flex flex-col gap-6 items-start justify-start w-full">
              <Typography variant="heading-xxlarge" color="dark">
                Typography System
              </Typography>
              <div className="bg-[#e0e0e0] h-0.5 w-full" />
            </div>
          </div>
          
          {/* Regular Headings */}
          <div className="w-full">
            <div className="flex flex-col gap-8 items-start justify-start w-full">
              <div className="flex flex-row items-start justify-between w-full">
                <div className="flex flex-col gap-4">
                  <Typography variant="heading-xsmall" color="light">xsmall</Typography>
                  <Typography variant="heading-xsmall" color="dark">
                    The quick brown fox jumps over the lazy dog
                  </Typography>
                </div>
                <div className="w-[280px]">
                  <Typography variant="heading-xsmall" color="dark" className="mb-2">
                    6 heading sizes
                  </Typography>
                  <Typography variant="text-large-paragraph" color="light">
                    SF Pro Display is used after 21px, which includes the default size and up.
                  </Typography>
                </div>
              </div>
              
              <div className="flex flex-row items-start justify-between w-full">
                <div className="flex flex-col gap-4">
                  <Typography variant="heading-xsmall" color="light">small</Typography>
                  <Typography variant="heading-small" color="dark">
                    The quick brown fox jumps over the lazy dog
                  </Typography>
                </div>
              </div>
              
              <div className="flex flex-row items-start justify-between w-full">
                <div className="flex flex-col gap-4">
                  <Typography variant="heading-xsmall" color="red">default</Typography>
                  <Typography variant="heading-default" color="dark">
                    The quick brown fox jumps over the lazy dog
                  </Typography>
                </div>
              </div>
              
              <div className="flex flex-row items-start justify-between w-full">
                <div className="flex flex-col gap-4">
                  <Typography variant="heading-xsmall" color="light">large</Typography>
                  <Typography variant="heading-large" color="dark">
                    The quick brown fox jumps over the lazy dog
                  </Typography>
                </div>
              </div>
              
              <div className="flex flex-row items-start justify-between w-full">
                <div className="flex flex-col gap-4">
                  <Typography variant="heading-xsmall" color="light">xlarge</Typography>
                  <Typography variant="heading-xlarge" color="dark">
                    The quick brown fox jumps over the lazy dog
                  </Typography>
                </div>
              </div>
              
              <div className="flex flex-row items-start justify-between w-full">
                <div className="flex flex-col gap-4">
                  <Typography variant="heading-xsmall" color="light">xxlarge</Typography>
                  <Typography variant="heading-xxlarge" color="dark">
                    The quick brown fox jumps over the lazy dog
                  </Typography>
                </div>
              </div>
            </div>
          </div>

          {/* All Caps Headings */}
          <div className="w-full">
            <div className="flex flex-col gap-8 items-start justify-start w-full">
              <div className="flex flex-row items-start justify-between w-full">
                <div className="flex flex-col gap-4">
                  <Typography variant="heading-xsmall" color="light">xsmall</Typography>
                  <Typography variant="heading-xsmall-caps" color="dark">
                    The quick brown fox jumps over the lazy dog
                  </Typography>
                </div>
                <div className="w-[280px]">
                  <Typography variant="heading-xsmall" color="dark" className="mb-2">
                    <span>3 </span>
                    <span className="text-[#ef3061]">all caps</span>
                    <span> heading sizes</span>
                  </Typography>
                  <Typography variant="text-large-paragraph" color="light">
                    SF Pro Text is used for all sizes since they are smaller than 21px.
                  </Typography>
                </div>
              </div>
              
              <div className="flex flex-row items-start justify-between w-full">
                <div className="flex flex-col gap-4">
                  <Typography variant="heading-xsmall" color="light">small</Typography>
                  <Typography variant="heading-small-caps" color="dark">
                    The quick brown fox jumps over the lazy dog
                  </Typography>
                </div>
              </div>
              
              <div className="flex flex-row items-start justify-between w-full">
                <div className="flex flex-col gap-4">
                  <Typography variant="heading-xsmall" color="red">default</Typography>
                  <Typography variant="heading-default-caps" color="dark">
                    The quick brown fox jumps over the lazy dog
                  </Typography>
                </div>
              </div>
            </div>
          </div>

          {/* Text Sizes */}
          <div className="w-full">
            <div className="flex flex-col gap-8 items-start justify-start w-full">
              <Typography variant="heading-default" color="dark">Text Sizes</Typography>
              
              {/* Single Line Text Examples */}
              <div className="flex flex-col gap-4">
                <Typography variant="heading-small" color="dark">Single Line Text</Typography>
                <div className="flex flex-col gap-2">
                  <Typography variant="text-small" color="dark">
                    Text Small (11px) - The quick brown fox jumps over the lazy dog
                  </Typography>
                  <Typography variant="text-default" color="red">
                    Text Default (13px) - The quick brown fox jumps over the lazy dog
                  </Typography>
                  <Typography variant="text-large" color="dark">
                    Text Large (15px) - The quick brown fox jumps over the lazy dog
                  </Typography>
                  <Typography variant="text-xlarge" color="dark">
                    Text XLarge (17px) - The quick brown fox jumps over the lazy dog
                  </Typography>
                </div>
              </div>

              {/* Paragraph Text Examples */}
              <div className="flex flex-col gap-4">
                <Typography variant="heading-small" color="dark">Paragraph Text</Typography>
                <div className="flex flex-col gap-4 max-w-[640px]">
                  <Typography variant="text-small-paragraph" color="dark">
                    Text Small Paragraph (11px) - Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
                  </Typography>
                  <Typography variant="text-default-paragraph" color="red">
                    Text Default Paragraph (13px) - Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
                  </Typography>
                  <Typography variant="text-large-paragraph" color="dark">
                    Text Large Paragraph (15px) - Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
                  </Typography>
                  <Typography variant="text-xlarge-paragraph" color="dark">
                    Text XLarge Paragraph (17px) - Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
                  </Typography>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}; 