import React from 'react';
import { AppColor, getColorClass } from './Colors';
import { Button, EditIcon, Tooltip } from './Button';

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
  color?: AppColor;
  className?: string;
}

export const Typography: React.FC<TypographyProps> = ({ 
  children, 
  variant = 'text-default', 
  color = 'dark',
  className = '' 
}) => {
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
    <div className={`${variantClasses[variant]} ${getColorClass(color)} ${className}`}>
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

          {/* Color Examples */}
          <div className="w-full">
            <div className="flex flex-col gap-8 items-start justify-start w-full">
              <Typography variant="heading-default" color="dark">Color Examples</Typography>
              
              {/* Grayscale Colors */}
              <div className="flex flex-col gap-4">
                <Typography variant="heading-small" color="dark">Grayscale Colors</Typography>
                <div className="flex flex-col gap-2">
                  <Typography variant="text-large" color="dark">
                    Dark text (#333333)
                  </Typography>
                  <Typography variant="text-large" color="light">
                    Light text (#757575)
                  </Typography>
                  <Typography variant="text-large" color="light-gray-4">
                    Light Gray 4 text (#e0e0e0)
                  </Typography>
                  <div className="bg-dark p-2 rounded">
                    <Typography variant="text-large" color="white">
                      White text on dark background (#ffffff)
                    </Typography>
                  </div>
                </div>
              </div>

              {/* Primary Color Families */}
              <div className="flex flex-col gap-6">
                <Typography variant="heading-small" color="dark">Primary Color Families</Typography>
                
                {/* Blue Colors */}
                <div className="flex flex-col gap-2">
                  <Typography variant="heading-xsmall" color="blue">Blue Family</Typography>
                  <div className="flex flex-col gap-1">
                    <Typography variant="text-default" color="blue-bright">
                      Blue Bright - Interactive elements and links
                    </Typography>
                    <Typography variant="text-default" color="blue">
                      Blue Default - Primary actions
                    </Typography>
                    <Typography variant="text-default" color="blue-dark-1">
                      Blue Dark - Hover states
                    </Typography>
                  </div>
                </div>

                {/* Red Colors */}
                <div className="flex flex-col gap-2">
                  <Typography variant="heading-xsmall" color="red">Red Family</Typography>
                  <div className="flex flex-col gap-1">
                    <Typography variant="text-default" color="red-bright">
                      Red Bright - Critical alerts and errors
                    </Typography>
                    <Typography variant="text-default" color="red">
                      Red Default - Error messages
                    </Typography>
                    <Typography variant="text-default" color="red-dark-1">
                      Red Dark - Error hover states
                    </Typography>
                  </div>
                </div>

                {/* Green Colors */}
                <div className="flex flex-col gap-2">
                  <Typography variant="heading-xsmall" color="green">Green Family</Typography>
                  <div className="flex flex-col gap-1">
                    <Typography variant="text-default" color="green-bright">
                      Green Bright - Success notifications
                    </Typography>
                    <Typography variant="text-default" color="green">
                      Green Default - Positive actions
                    </Typography>
                    <Typography variant="text-default" color="green-dark-1">
                      Green Dark - Success hover states
                    </Typography>
                  </div>
                </div>

                {/* Orange Colors */}
                <div className="flex flex-col gap-2">
                  <Typography variant="heading-xsmall" color="orange">Orange Family</Typography>
                  <div className="flex flex-col gap-1">
                    <Typography variant="text-default" color="orange-bright">
                      Orange Bright - Warning alerts
                    </Typography>
                    <Typography variant="text-default" color="orange">
                      Orange Default - Warning messages
                    </Typography>
                    <Typography variant="text-default" color="orange-dark-1">
                      Orange Dark - Warning hover states
                    </Typography>
                  </div>
                </div>
              </div>

              {/* Accent Colors */}
              <div className="flex flex-col gap-4">
                <Typography variant="heading-small" color="dark">Accent Colors</Typography>
                <div className="grid grid-cols-2 gap-6">
                  <div className="flex flex-col gap-2">
                    <Typography variant="text-default" color="purple">
                      Purple - Premium features
                    </Typography>
                    <Typography variant="text-default" color="pink">
                      Pink - Special highlights
                    </Typography>
                    <Typography variant="text-default" color="cyan">
                      Cyan - Information states
                    </Typography>
                  </div>
                  <div className="flex flex-col gap-2">
                    <Typography variant="text-default" color="teal">
                      Teal - Secondary actions
                    </Typography>
                    <Typography variant="text-default" color="yellow">
                      Yellow - Attention markers
                    </Typography>
                  </div>
                </div>
              </div>

              {/* Heading Color Examples */}
              <div className="flex flex-col gap-4">
                <Typography variant="heading-small" color="dark">Heading Color Examples</Typography>
                <div className="flex flex-col gap-3">
                  <Typography variant="heading-large" color="blue-bright">
                    Blue Bright Heading
                  </Typography>
                  <Typography variant="heading-default" color="purple">
                    Purple Default Heading
                  </Typography>
                  <Typography variant="heading-small" color="green">
                    Green Small Heading
                  </Typography>
                  <Typography variant="heading-xsmall-caps" color="red">
                    Red All Caps Heading
                  </Typography>
                </div>
              </div>
            </div>
          </div>

          {/* Button Examples */}
          <div className="w-full">
            <div className="flex flex-col gap-10 items-start justify-start w-full">
              <div className="w-full">
                <div className="flex flex-col gap-6 items-start justify-start w-full">
                  <Typography variant="heading-xxlarge" color="dark">
                    Buttons
                  </Typography>
                  <div className="bg-[#e0e0e0] h-0.5 w-full" />
                </div>
              </div>

              <div className="flex flex-row gap-16 items-start justify-start">
                {/* Regular Buttons Column */}
                <div className="flex flex-col gap-16">
                  {/* Button / Default */}
                  <div className="flex flex-col gap-4">
                    <Typography variant="heading-xsmall-caps" color="light">
                      BUTTON / DEFAULT
                    </Typography>
                    <div className="flex flex-col gap-4">
                      <div className="flex flex-row gap-4 items-center">
                        <Button variant="default">Default</Button>
                        <Button variant="primary">Primary</Button>
                        <Button variant="danger">Danger</Button>
                        <Button variant="secondary">Secondary</Button>
                      </div>
                      <div className="flex flex-row gap-4 items-center">
                        <Button variant="default" icon={<EditIcon />}>Default</Button>
                        <Button variant="primary" icon={<EditIcon />}>Primary</Button>
                        <Button variant="danger" icon={<EditIcon />}>Danger</Button>
                        <Button variant="secondary" icon={<EditIcon />}>Secondary</Button>
                      </div>
                    </div>
                  </div>

                  {/* Button / Small */}
                  <div className="flex flex-col gap-4">
                    <Typography variant="heading-xsmall-caps" color="light">
                      BUTTON / SMALL
                    </Typography>
                    <div className="flex flex-col gap-4">
                      <div className="flex flex-row gap-4 items-center">
                        <Button variant="default" size="small">Default</Button>
                        <Button variant="primary" size="small">Primary</Button>
                        <Button variant="danger" size="small">Danger</Button>
                        <Button variant="secondary" size="small">Secondary</Button>
                      </div>
                      <div className="flex flex-row gap-4 items-center">
                        <Button variant="default" size="small" icon={<EditIcon />}>Default</Button>
                        <Button variant="primary" size="small" icon={<EditIcon />}>Primary</Button>
                        <Button variant="danger" size="small" icon={<EditIcon />}>Danger</Button>
                        <Button variant="secondary" size="small" icon={<EditIcon />}>Secondary</Button>
                      </div>
                    </div>
                  </div>

                  {/* Button / Large */}
                  <div className="flex flex-col gap-4">
                    <Typography variant="heading-xsmall-caps" color="light">
                      BUTTON / LARGE
                    </Typography>
                    <div className="flex flex-col gap-4">
                      <div className="flex flex-row gap-4 items-center">
                        <Button variant="default" size="large">Default</Button>
                        <Button variant="primary" size="large">Primary</Button>
                        <Button variant="danger" size="large">Danger</Button>
                        <Button variant="secondary" size="large">Secondary</Button>
                      </div>
                      <div className="flex flex-row gap-4 items-center">
                        <Button variant="default" size="large" icon={<EditIcon />}>Default</Button>
                        <Button variant="primary" size="large" icon={<EditIcon />}>Primary</Button>
                        <Button variant="danger" size="large" icon={<EditIcon />}>Danger</Button>
                        <Button variant="secondary" size="large" icon={<EditIcon />}>Secondary</Button>
                      </div>
                    </div>
                  </div>

                  {/* Icon-only button with tooltip */}
                  <div className="flex flex-col gap-4">
                    <Typography variant="heading-xsmall-caps" color="light">
                      ICON-ONLY BUTTON WITH TOOLTIP
                    </Typography>
                    <div className="flex flex-row gap-4 items-center">
                      <Tooltip text="Edit content">
                        <Button variant="default" iconOnly icon={<EditIcon />}>
                          Edit
                        </Button>
                      </Tooltip>
                    </div>
                  </div>
                </div>

                {/* Text Buttons Column */}
                <div className="flex flex-col gap-16">
                  <div className="w-[323px]">
                    <Typography variant="text-large-paragraph" color="light">
                      Text buttons can be also be used for links.
                    </Typography>
                  </div>

                  {/* Text Button / Default */}
                  <div className="flex flex-col gap-4">
                    <Typography variant="heading-xsmall-caps" color="light">
                      TEXT BUTTON / DEFAULT
                    </Typography>
                    <div className="flex flex-col gap-2">
                      <Button variant="text" size="default">Default</Button>
                      <Button variant="text" size="default" className="text-[#333333] hover:text-[#757575]">Dark</Button>
                      <Button variant="text" size="default" className="text-[#757575] hover:text-[#333333]">Light</Button>
                    </div>
                    <div className="flex flex-col gap-2">
                      <Button variant="text" size="default" icon={<EditIcon />}>Default</Button>
                      <Button variant="text" size="default" icon={<EditIcon />} className="text-[#333333] hover:text-[#757575]">Dark</Button>
                      <Button variant="text" size="default" icon={<EditIcon />} className="text-[#757575] hover:text-[#333333]">Light</Button>
                    </div>
                  </div>

                  {/* Text Button / Small */}
                  <div className="flex flex-col gap-4">
                    <Typography variant="heading-xsmall-caps" color="light">
                      TEXT BUTTON / SMALL
                    </Typography>
                    <div className="flex flex-col gap-2">
                      <Button variant="text" size="small">Default</Button>
                      <Button variant="text" size="small" className="text-[#333333] hover:text-[#757575]">Dark</Button>
                      <Button variant="text" size="small" className="text-[#757575] hover:text-[#333333]">Light</Button>
                    </div>
                    <div className="flex flex-col gap-2">
                      <Button variant="text" size="small" icon={<EditIcon />}>Default</Button>
                      <Button variant="text" size="small" icon={<EditIcon />} className="text-[#333333] hover:text-[#757575]">Dark</Button>
                      <Button variant="text" size="small" icon={<EditIcon />} className="text-[#757575] hover:text-[#333333]">Light</Button>
                    </div>
                  </div>

                  {/* Text Button / Large */}
                  <div className="flex flex-col gap-4">
                    <Typography variant="heading-xsmall-caps" color="light">
                      TEXT BUTTON / LARGE
                    </Typography>
                    <div className="flex flex-col gap-2">
                      <Button variant="text" size="large">Default</Button>
                      <Button variant="text" size="large" className="text-[#333333] hover:text-[#757575]">Dark</Button>
                      <Button variant="text" size="large" className="text-[#757575] hover:text-[#333333]">Light</Button>
                    </div>
                    <div className="flex flex-col gap-2">
                      <Button variant="text" size="large" icon={<EditIcon />}>Default</Button>
                      <Button variant="text" size="large" icon={<EditIcon />} className="text-[#333333] hover:text-[#757575]">Dark</Button>
                      <Button variant="text" size="large" icon={<EditIcon />} className="text-[#757575] hover:text-[#333333]">Light</Button>
                    </div>
                  </div>

                  {/* Text Button / XLarge */}
                  <div className="flex flex-col gap-4">
                    <Typography variant="heading-xsmall-caps" color="light">
                      TEXT BUTTON / XLARGE
                    </Typography>
                    <div className="flex flex-col gap-2">
                      <Button variant="text" size="xlarge">Default</Button>
                      <Button variant="text" size="xlarge" className="text-[#333333] hover:text-[#757575]">Dark</Button>
                      <Button variant="text" size="xlarge" className="text-[#757575] hover:text-[#333333]">Light</Button>
                    </div>
                    <div className="flex flex-col gap-2">
                      <Button variant="text" size="xlarge" icon={<EditIcon />}>Default</Button>
                      <Button variant="text" size="xlarge" icon={<EditIcon />} className="text-[#333333] hover:text-[#757575]">Dark</Button>
                      <Button variant="text" size="xlarge" icon={<EditIcon />} className="text-[#757575] hover:text-[#333333]">Light</Button>
                    </div>
                  </div>

                  {/* Icon-only text button with tooltip */}
                  <div className="flex flex-col gap-4">
                    <Typography variant="heading-xsmall-caps" color="light">
                      ICON-ONLY TEXT BUTTON WITH TOOLTIP
                    </Typography>
                    <div className="flex flex-row gap-4 items-center">
                      <Tooltip text="Edit content">
                        <Button variant="text" iconOnly icon={<EditIcon />}>
                          Edit
                        </Button>
                      </Tooltip>
                    </div>
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