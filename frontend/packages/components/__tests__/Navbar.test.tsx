import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { useSession } from 'next-auth/react';
import Navbar from '../Navbar';

jest.mock('next-auth/react');

const mockUseSession = useSession as jest.MockedFunction<typeof useSession>;

describe('Navbar Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    (global.fetch as jest.Mock).mockClear();
  });

  it('should match snapshot when not authenticated', () => {
    mockUseSession.mockReturnValue({
      data: null,
      status: 'unauthenticated',
      update: jest.fn(),
    });

    const { container } = render(<Navbar />);
    expect(container).toMatchSnapshot();
  });

  it('should display app title', () => {
    mockUseSession.mockReturnValue({
      data: null,
      status: 'unauthenticated',
      update: jest.fn(),
    });

    render(<Navbar />);
    expect(screen.getByText('EVE Industry Tool')).toBeInTheDocument();
  });

  it('should have all navigation links', () => {
    mockUseSession.mockReturnValue({
      data: null,
      status: 'unauthenticated',
      update: jest.fn(),
    });

    render(<Navbar />);

    const charactersLink = screen.getByRole('link', { name: /characters/i });
    const corporationsLink = screen.getByRole('link', { name: /corporations/i });
    const inventoryLink = screen.getByRole('link', { name: /inventory/i });
    const stockpilesLink = screen.getByRole('link', { name: /stockpiles/i });
    const contactsLink = screen.getByRole('link', { name: /contacts/i });

    expect(charactersLink).toHaveAttribute('href', '/characters');
    expect(corporationsLink).toHaveAttribute('href', '/corporations');
    expect(inventoryLink).toHaveAttribute('href', '/inventory');
    expect(stockpilesLink).toHaveAttribute('href', '/stockpiles');
    expect(contactsLink).toHaveAttribute('href', '/contacts');
  });

  it('should have rocket icon', () => {
    mockUseSession.mockReturnValue({
      data: null,
      status: 'unauthenticated',
      update: jest.fn(),
    });

    render(<Navbar />);
    const menuButton = screen.getByLabelText('menu');
    expect(menuButton).toBeInTheDocument();
  });

  it('should not fetch contacts when user is not authenticated', () => {
    mockUseSession.mockReturnValue({
      data: null,
      status: 'unauthenticated',
      update: jest.fn(),
    });

    render(<Navbar />);
    expect(global.fetch).not.toHaveBeenCalled();
  });

  it('should fetch contacts when user is authenticated', async () => {
    mockUseSession.mockReturnValue({
      data: { providerAccountId: '123' } as any,
      status: 'authenticated',
      update: jest.fn(),
    });

    (global.fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: async () => [],
    });

    render(<Navbar />);

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledWith('/api/contacts');
    });
  });

  it('should display badge with pending contact count', async () => {
    mockUseSession.mockReturnValue({
      data: { providerAccountId: '123' } as any,
      status: 'authenticated',
      update: jest.fn(),
    });

    const mockContacts = [
      {
        id: 1,
        requesterUserId: 456,
        recipientUserId: 123,
        requesterName: 'Test User 1',
        recipientName: 'Current User',
        status: 'pending',
        requestedAt: '2024-01-01',
      },
      {
        id: 2,
        requesterUserId: 789,
        recipientUserId: 123,
        requesterName: 'Test User 2',
        recipientName: 'Current User',
        status: 'pending',
        requestedAt: '2024-01-02',
      },
      {
        id: 3,
        requesterUserId: 123,
        recipientUserId: 999,
        requesterName: 'Current User',
        recipientName: 'Test User 3',
        status: 'pending',
        requestedAt: '2024-01-03',
      },
    ];

    (global.fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: async () => mockContacts,
    });

    render(<Navbar />);

    await waitFor(() => {
      const badge = screen.getByText('2');
      expect(badge).toBeInTheDocument();
    });
  });

  it('should not display badge when there are no pending requests', async () => {
    mockUseSession.mockReturnValue({
      data: { providerAccountId: '123' } as any,
      status: 'authenticated',
      update: jest.fn(),
    });

    const mockContacts = [
      {
        id: 1,
        requesterUserId: 123,
        recipientUserId: 456,
        requesterName: 'Current User',
        recipientName: 'Test User',
        status: 'accepted',
        requestedAt: '2024-01-01',
      },
    ];

    (global.fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: async () => mockContacts,
    });

    render(<Navbar />);

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalled();
    });

    // Badge should not be visible when count is 0
    expect(screen.queryByText('0')).not.toBeInTheDocument();
  });

  it('should handle fetch errors gracefully', async () => {
    mockUseSession.mockReturnValue({
      data: { providerAccountId: '123' } as any,
      status: 'authenticated',
      update: jest.fn(),
    });

    (global.fetch as jest.Mock).mockResolvedValue({
      ok: false,
    });

    const consoleErrorSpy = jest.spyOn(console, 'error').mockImplementation();

    render(<Navbar />);

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalled();
    });

    // Should not crash on error
    expect(screen.getByText('EVE Industry Tool')).toBeInTheDocument();

    consoleErrorSpy.mockRestore();
  });

  it('should set up polling interval for contact updates', async () => {
    jest.useFakeTimers();

    mockUseSession.mockReturnValue({
      data: { providerAccountId: '123' } as any,
      status: 'authenticated',
      update: jest.fn(),
    });

    (global.fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: async () => [],
    });

    render(<Navbar />);

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledTimes(1);
    });

    // Fast forward 30 seconds
    jest.advanceTimersByTime(30000);

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledTimes(2);
    });

    // Fast forward another 30 seconds
    jest.advanceTimersByTime(30000);

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledTimes(3);
    });

    jest.useRealTimers();
  });
});
