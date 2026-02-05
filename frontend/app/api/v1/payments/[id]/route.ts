import { NextRequest, NextResponse } from 'next/server'

const PAYMENT_SERVICE_URL = process.env.PAYMENT_SERVICE_URL || 'http://localhost:8083'

export async function GET(
    request: NextRequest,
    { params }: { params: Promise<{ id: string }> }
) {
    try {
        const { id } = await params

        const response = await fetch(`${PAYMENT_SERVICE_URL}/api/v1/protected/payments/${id}`, {
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
                { error: 'Failed to fetch payment details', details: errorData },
                { status: response.status }
            )
        }

        const data = await response.json()
        const result = data?.data || data
        return NextResponse.json(result)
    } catch (error) {
        console.error('Get payment details API error:', error)
        return NextResponse.json(
            { error: 'Internal server error' },
            { status: 500 }
        )
    }
}

export async function PUT(
    request: NextRequest,
    { params }: { params: Promise<{ id: string }> }
) {
    try {
        const { id } = await params
        const body = await request.json()

        // Check if it's a status update
        const url = request.url.endsWith('/status')
            ? `${PAYMENT_SERVICE_URL}/api/v1/protected/payments/${id}/status`
            : `${PAYMENT_SERVICE_URL}/api/v1/protected/payments/${id}`

        const response = await fetch(url, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Cookie': request.headers.get('cookie') || '',
                'Authorization': request.headers.get('authorization') || '',
            },
            body: JSON.stringify(body),
        })

        if (!response.ok) {
            const errorData = await response.text()
            return NextResponse.json(
                { error: 'Failed to update payment', details: errorData },
                { status: response.status }
            )
        }

        const data = await response.json()
        const result = data?.data || data
        return NextResponse.json(result)
    } catch (error) {
        console.error('Update payment API error:', error)
        return NextResponse.json(
            { error: 'Internal server error' },
            { status: 500 }
        )
    }
}
