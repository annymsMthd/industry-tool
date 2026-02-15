import React, { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Button,
  ButtonGroup,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  CircularProgress,
  Alert,
  Stack,
} from '@mui/material';
import DownloadIcon from '@mui/icons-material/Download';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import AttachMoneyIcon from '@mui/icons-material/AttachMoney';
import CategoryIcon from '@mui/icons-material/Category';
import PeopleIcon from '@mui/icons-material/People';

interface TimeSeriesData {
  date: string;
  revenue: number;
  transactions: number;
  quantitySold: number;
}

interface ItemSalesData {
  typeId: number;
  typeName: string;
  quantitySold: number;
  revenue: number;
  transactionCount: number;
  averagePricePerUnit: number;
}

interface SalesMetrics {
  totalRevenue: number;
  totalTransactions: number;
  totalQuantitySold: number;
  uniqueItemTypes: number;
  uniqueBuyers: number;
  timeSeriesData: TimeSeriesData[];
  topItems: ItemSalesData[];
}

type TimePeriod = '7d' | '30d' | '90d' | '1y' | 'all';

const formatISK = (value: number): string => {
  if (value >= 1_000_000_000) {
    return `${(value / 1_000_000_000).toFixed(2)}B ISK`;
  } else if (value >= 1_000_000) {
    return `${(value / 1_000_000).toFixed(2)}M ISK`;
  } else if (value >= 1_000) {
    return `${(value / 1_000).toFixed(2)}K ISK`;
  }
  return `${value.toFixed(2)} ISK`;
};

const formatNumber = (value: number): string => {
  return value.toLocaleString();
};

export default function SalesMetrics() {
  const [period, setPeriod] = useState<TimePeriod>('30d');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [metrics, setMetrics] = useState<SalesMetrics | null>(null);

  useEffect(() => {
    fetchMetrics();
  }, [period]);

  const fetchMetrics = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`/api/analytics/sales?period=${period}`);
      if (!response.ok) {
        throw new Error('Failed to fetch sales metrics');
      }
      const data = await response.json();
      setMetrics(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const exportToCSV = () => {
    if (!metrics) return;

    // Export time series data
    const csvRows = [
      ['Date', 'Revenue (ISK)', 'Transactions', 'Quantity Sold'],
      ...metrics.timeSeriesData.map(row => [
        row.date,
        row.revenue.toString(),
        row.transactions.toString(),
        row.quantitySold.toString(),
      ]),
    ];

    const csvContent = csvRows.map(row => row.join(',')).join('\n');
    const blob = new Blob([csvContent], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `sales-metrics-${period}-${new Date().toISOString().split('T')[0]}.csv`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Alert severity="error" sx={{ mb: 2 }}>
        {error}
      </Alert>
    );
  }

  if (!metrics) {
    return (
      <Alert severity="info" sx={{ mb: 2 }}>
        No sales data available
      </Alert>
    );
  }

  return (
    <Box>
      {/* Header with time period filter and export button */}
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Sales Analytics</Typography>
        <Box display="flex" gap={2}>
          <ButtonGroup variant="outlined" size="small">
            {(['7d', '30d', '90d', '1y', 'all'] as TimePeriod[]).map((p) => (
              <Button
                key={p}
                onClick={() => setPeriod(p)}
                variant={period === p ? 'contained' : 'outlined'}
              >
                {p === 'all' ? 'All Time' : p.toUpperCase()}
              </Button>
            ))}
          </ButtonGroup>
          <Button
            variant="outlined"
            startIcon={<DownloadIcon />}
            onClick={exportToCSV}
          >
            Export CSV
          </Button>
        </Box>
      </Box>

      {/* Info Alert */}
      <Alert severity="info" sx={{ mb: 3 }}>
        Analytics only shows completed transactions. Pending sales (contract_created status) are not included in these metrics.
      </Alert>

      {/* Summary Cards */}
      <Box
        sx={{
          display: 'grid',
          gridTemplateColumns: { xs: '1fr', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)' },
          gap: 3,
          mb: 4,
        }}
      >
        <Box>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <AttachMoneyIcon color="primary" sx={{ mr: 1 }} />
                <Typography color="textSecondary" variant="body2">
                  Total Revenue
                </Typography>
              </Box>
              <Typography variant="h5">{formatISK(metrics.totalRevenue)}</Typography>
            </CardContent>
          </Card>
        </Box>

        <Box>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <ShoppingCartIcon color="primary" sx={{ mr: 1 }} />
                <Typography color="textSecondary" variant="body2">
                  Total Transactions
                </Typography>
              </Box>
              <Typography variant="h5">{formatNumber(metrics.totalTransactions)}</Typography>
            </CardContent>
          </Card>
        </Box>

        <Box>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <TrendingUpIcon color="primary" sx={{ mr: 1 }} />
                <Typography color="textSecondary" variant="body2">
                  Items Sold
                </Typography>
              </Box>
              <Typography variant="h5">{formatNumber(metrics.totalQuantitySold)}</Typography>
            </CardContent>
          </Card>
        </Box>

        <Box>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <CategoryIcon color="primary" sx={{ mr: 1 }} />
                <Typography color="textSecondary" variant="body2">
                  Unique Item Types
                </Typography>
              </Box>
              <Typography variant="h5">{formatNumber(metrics.uniqueItemTypes)}</Typography>
            </CardContent>
          </Card>
        </Box>

        <Box>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <PeopleIcon color="primary" sx={{ mr: 1 }} />
                <Typography color="textSecondary" variant="body2">
                  Unique Buyers
                </Typography>
              </Box>
              <Typography variant="h5">{formatNumber(metrics.uniqueBuyers)}</Typography>
            </CardContent>
          </Card>
        </Box>

        <Box>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <AttachMoneyIcon color="primary" sx={{ mr: 1 }} />
                <Typography color="textSecondary" variant="body2">
                  Avg Revenue/Transaction
                </Typography>
              </Box>
              <Typography variant="h5">
                {metrics.totalTransactions > 0
                  ? formatISK(metrics.totalRevenue / metrics.totalTransactions)
                  : '0 ISK'}
              </Typography>
            </CardContent>
          </Card>
        </Box>
      </Box>

      {/* Top Selling Items */}
      <Card sx={{ mb: 4 }}>
        <CardContent>
          <Typography variant="h6" mb={2}>
            Top Selling Items
          </Typography>
          <TableContainer component={Paper} variant="outlined">
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Item Name</TableCell>
                  <TableCell align="right">Quantity Sold</TableCell>
                  <TableCell align="right">Revenue</TableCell>
                  <TableCell align="right">Transactions</TableCell>
                  <TableCell align="right">Avg Price/Unit</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {metrics.topItems.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={5} align="center">
                      No sales data available
                    </TableCell>
                  </TableRow>
                ) : (
                  metrics.topItems.map((item) => (
                    <TableRow key={item.typeId}>
                      <TableCell>{item.typeName}</TableCell>
                      <TableCell align="right">{formatNumber(item.quantitySold)}</TableCell>
                      <TableCell align="right">{formatISK(item.revenue)}</TableCell>
                      <TableCell align="right">{formatNumber(item.transactionCount)}</TableCell>
                      <TableCell align="right">{formatISK(item.averagePricePerUnit)}</TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </TableContainer>
        </CardContent>
      </Card>

      {/* Sales Over Time */}
      <Card>
        <CardContent>
          <Typography variant="h6" mb={2}>
            Sales Over Time
          </Typography>
          <TableContainer component={Paper} variant="outlined">
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Date</TableCell>
                  <TableCell align="right">Revenue</TableCell>
                  <TableCell align="right">Transactions</TableCell>
                  <TableCell align="right">Quantity Sold</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {metrics.timeSeriesData.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={4} align="center">
                      No time series data available
                    </TableCell>
                  </TableRow>
                ) : (
                  metrics.timeSeriesData.map((row) => (
                    <TableRow key={row.date}>
                      <TableCell>{row.date}</TableCell>
                      <TableCell align="right">{formatISK(row.revenue)}</TableCell>
                      <TableCell align="right">{formatNumber(row.transactions)}</TableCell>
                      <TableCell align="right">{formatNumber(row.quantitySold)}</TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </TableContainer>
        </CardContent>
      </Card>
    </Box>
  );
}
