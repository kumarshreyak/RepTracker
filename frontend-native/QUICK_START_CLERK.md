# Quick Start: Clerk Authentication

## 🚀 Get Started in 5 Minutes

### Step 1: Get Your Clerk Key (2 minutes)
1. Go to **https://clerk.com** → Sign up/Sign in
2. Click **"Add application"**
3. Name it **"GymLog"** → Select **"Email"** → Click **"Create"**
4. Go to **"API Keys"** → Copy the **Publishable Key** (starts with `pk_test_`)

### Step 2: Configure Environment (1 minute)
Create a `.env` file in `frontend-native/` directory:

```bash
EXPO_PUBLIC_API_BASE_URL=http://localhost:8080
EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_YOUR_KEY_HERE
```

**Replace `pk_test_YOUR_KEY_HERE` with your actual key!**

### Step 3: Start the App (2 minutes)
```bash
cd frontend-native
npm start
```

### Step 4: Test It Out
1. **Sign Up**: Create a new account with your email
2. **Check Email**: Look for verification code (check spam!)
3. **Verify**: Enter the 6-digit code
4. **Sign In**: Use your credentials to sign in
5. **Sign Out**: Go to Profile tab → Sign Out

## ✅ That's It!

Your app now has secure authentication with:
- ✅ Email/password sign-up
- ✅ Email verification
- ✅ Secure session management
- ✅ Protected routes
- ✅ Sign-out functionality

## 🆘 Quick Troubleshooting

**"Missing EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY"**
→ Check your `.env` file exists and has the correct key

**"Email not received"**
→ Check spam folder, verify email settings in Clerk dashboard

**"Can't sign in"**
→ Make sure you completed email verification during sign-up

## 📖 Need More Help?

- **Full Setup Guide**: See `CLERK_SETUP.md`
- **Migration Details**: See `CLERK_MIGRATION_SUMMARY.md`
- **Clerk Docs**: https://clerk.com/docs/quickstarts/expo

