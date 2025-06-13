import { Typography, LoginButton } from "@/components";

export default function Home() {
  return (
    <main className="min-h-screen bg-light-gray-1 flex items-center justify-center">
      <div className="max-w-md w-full mx-auto px-6 py-8">
        <div className="text-center space-y-8">
          {/* Logo and Title */}
          <div className="space-y-3">
            <img 
              src="/gymlog-logo.svg" 
              alt="GymLog" 
              className="mx-auto h-16 w-auto"
            />
            <Typography variant="text-large-paragraph" color="light" className="max-w-sm mx-auto">
              Where every rep is counted
            </Typography>
          </div>

          {/* Login Button */}
          <div className="space-y-4">
            <LoginButton />
            <Typography variant="text-small" color="light" className="text-center">
              By signing in, you agree to our{" "}
              <a 
                href="https://docs.google.com/document/d/1khuoDSOc10TMlJ12M7heMVlvv6XeKKSdEpES1yvAwos/edit?usp=sharing" 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-blue-bright hover:text-blue underline"
              >
                terms of service
              </a>
              {" "}and{" "}
              <a 
                href="https://docs.google.com/document/d/1khuoDSOc10TMlJ12M7heMVlvv6XeKKSdEpES1yvAwos/edit?usp=sharing" 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-blue-bright hover:text-blue underline"
              >
                privacy policy
              </a>
            </Typography>
          </div>
        </div>
      </div>
    </main>
  );
}
