import { NextRequest, NextResponse } from 'next/server'

export async function POST(req: NextRequest) {
  try {
    const { googleId, email, name, picture, accessToken, refreshToken, expiresAt } = await req.json()

    // Call backend HTTP API
    const backendUrl = process.env.BACKEND_HTTP_URL || 'http://localhost:8080'
    const response = await fetch(`${backendUrl}/api/auth/google`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        googleId,
        email,
        name,
        picture,
        accessToken,
        refreshToken,
        expiresAt,
      }),
    })

    if (!response.ok) {
      throw new Error(`Backend responded with status: ${response.status}`)
    }

    const user = await response.json()
    
    // Extract session token from response headers
    const sessionToken = response.headers.get('X-Session-Token')
    
    return NextResponse.json({
      ...user,
      sessionToken,
    })
  } catch (error) {
    console.error('Error in Google auth API:', error)
    return NextResponse.json(
      { error: 'Failed to authenticate with backend' },
      { status: 500 }
    )
  }
} 