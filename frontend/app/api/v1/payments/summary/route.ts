import { NextRequest, NextResponse } from 'next/server'

const PAYMENT_SERVICE_URL = process.env.PAYMENT_SERVICE_URL || 'http://localhost:8083'

export async function GET(request: NextRequest) {
    try {
        const { searchParams } = new URL(request.url)
        const queryString = searchParams.toString()

        const response = await fetch(`${PAYMENT_SERVICE_URL}/api/v1/protected/payments/summary${queryString ? `?${queryString}` : ''}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Cookie': request.headers.get('cookie') || '',
                'Authorization': request.headers.get('authorization') || '',
            },
        })

        if (!response.ok) {
            const errorData = await response.text()
            return NextResponse.json(
                { error: 'Failed to fetch payment summary', details: errorData },
                { status: response.status }
            )
        }

        const data = await response.json()
        const result = data?.data || data
        return NextResponse.json(result)
    } catch (error) {
        console.error('Payment summary API error:', error)
        return NextResponse.json(
            { error: 'Internal server error' },
            { status: 500 }
        )
    }
}
