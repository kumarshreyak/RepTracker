import { NextRequest, NextResponse } from 'next/server'
import { auth } from '@/auth'

export async function GET(req: NextRequest) {
  try {
    // Get session from NextAuth
    const session = await auth()
    
    if (!session?.user) {
      return NextResponse.json(
        { error: 'No session found' },
        { status: 401 }
      )
    }

    // Extract authorization header (session token from backend)
    const authHeader = req.headers.get('authorization')
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return NextResponse.json(
        { error: 'Missing or invalid authorization header' },
        { status: 401 }
      )
    }

    const token = authHeader.substring(7)
    
    // Validate session with backend
    const backendUrl = process.env.BACKEND_HTTP_URL || 'http://localhost:8080'
    const response = await fetch(`${backendUrl}/api/auth/validate`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Session expired or invalid' },
        { status: 401 }
      )
    }

    const user = await response.json()
    return NextResponse.json(user)
  } catch (error) {
    console.error('Error validating session:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
} 