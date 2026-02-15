import { useState, useEffect } from 'react';
import { useSession } from "next-auth/react";
import Loading from "@industry-tool/components/loading";
import Unauthorized from "@industry-tool/components/unauthorized";
import Navbar from "@industry-tool/components/Navbar";
import Container from '@mui/material/Container';
import Box from '@mui/material/Box';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import MyListings from "@industry-tool/components/marketplace/MyListings";
import MarketplaceBrowser from "@industry-tool/components/marketplace/MarketplaceBrowser";
import PurchaseHistory from "@industry-tool/components/marketplace/PurchaseHistory";
import PendingSales from "@industry-tool/components/marketplace/PendingSales";
import BuyOrders from "@industry-tool/components/marketplace/BuyOrders";
import DemandViewer from "@industry-tool/components/marketplace/DemandViewer";
import SalesMetrics from "@industry-tool/components/analytics/SalesMetrics";

export default function Marketplace() {
  const { status } = useSession();
  const [tabIndex, setTabIndex] = useState(() => {
    // Load saved tab from localStorage on initial render
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem('marketplaceTab');
      return saved ? parseInt(saved, 10) : 0;
    }
    return 0;
  });

  // Save tab selection to localStorage whenever it changes
  useEffect(() => {
    localStorage.setItem('marketplaceTab', tabIndex.toString());
  }, [tabIndex]);

  if (status === "loading") {
    return <Loading />;
  }

  if (status !== "authenticated") {
    return <Unauthorized />;
  }

  return (
    <>
      <Navbar />
      <Container maxWidth={false}>
        <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
          <Tabs value={tabIndex} onChange={(_, newValue) => setTabIndex(newValue)}>
            <Tab label="My Listings" />
            <Tab label="Browse" />
            <Tab label="Pending Sales" />
            <Tab label="History" />
            <Tab label="My Buy Orders" />
            <Tab label="Demand" />
            <Tab label="Analytics" />
          </Tabs>
        </Box>

        {tabIndex === 0 && <MyListings />}
        {tabIndex === 1 && <MarketplaceBrowser />}
        {tabIndex === 2 && <PendingSales />}
        {tabIndex === 3 && <PurchaseHistory />}
        {tabIndex === 4 && <BuyOrders />}
        {tabIndex === 5 && <DemandViewer />}
        {tabIndex === 6 && <SalesMetrics />}
      </Container>
    </>
  );
}
