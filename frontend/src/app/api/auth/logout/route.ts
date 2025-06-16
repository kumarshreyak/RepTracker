import { NextRequest, NextResponse } from 'next/server'

export async function POST(req: NextRequest) {
  try {
    // Extract authorization header (session token from backend)
    const authHeader = req.headers.get('authorization')
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return NextResponse.json(
        { error: 'Missing or invalid authorization header' },
        { status: 401 }
      )
    }

    const token = authHeader.substring(7)
    
    // Call backend to logout (delete session)
    const backendUrl = process.env.BACKEND_HTTP_URL || 'http://localhost:8080'
    const response = await fetch(`${backendUrl}/api/auth/logout`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (!response.ok) {
      return NextResponse.json(
        { error: 'Failed to logout from backend' },
        { status: 500 }
      )
    }

    return NextResponse.json({ message: 'Logged out successfully' })
  } catch (error) {
    console.error('Error during logout:', error)
    return NextResponse.json(
      { error: 'Internal server error' },
      { status: 500 }
    )
  }
} 