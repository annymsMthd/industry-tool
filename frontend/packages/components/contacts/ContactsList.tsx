import { useState, useEffect, useRef } from 'react';
import { useSession } from "next-auth/react";
import Navbar from "@industry-tool/components/Navbar";
import Loading from "@industry-tool/components/loading";
import Container from '@mui/material/Container';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Button from '@mui/material/Button';
import Snackbar from '@mui/material/Snackbar';
import Alert from '@mui/material/Alert';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Chip from '@mui/material/Chip';
import IconButton from '@mui/material/IconButton';
import DeleteIcon from '@mui/icons-material/Delete';
import CheckIcon from '@mui/icons-material/Check';
import CloseIcon from '@mui/icons-material/Close';
import SettingsIcon from '@mui/icons-material/Settings';
import PersonAddIcon from '@mui/icons-material/PersonAdd';
import Dialog from '@mui/material/Dialog';
import DialogTitle from '@mui/material/DialogTitle';
import DialogContent from '@mui/material/DialogContent';
import DialogActions from '@mui/material/DialogActions';
import TextField from '@mui/material/TextField';
import PermissionsDialog from './PermissionsDialog';

export type Contact = {
  id: number;
  requesterUserId: number;
  recipientUserId: number;
  requesterName: string;
  recipientName: string;
  status: 'pending' | 'accepted' | 'rejected';
  requestedAt: string;
  respondedAt?: string;
};

export default function ContactsList() {
  const { data: session } = useSession();
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [loading, setLoading] = useState(true);
  const [tabIndex, setTabIndex] = useState(0);
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState('');
  const [snackbarSeverity, setSnackbarSeverity] = useState<'success' | 'error'>('success');
  const [addContactOpen, setAddContactOpen] = useState(false);
  const [newContactCharacterName, setNewContactCharacterName] = useState('');
  const [permissionsDialogOpen, setPermissionsDialogOpen] = useState(false);
  const [selectedContact, setSelectedContact] = useState<Contact | null>(null);
  const hasFetchedRef = useRef(false);

  const currentUserId = session?.providerAccountId ? parseInt(session.providerAccountId) : null;

  useEffect(() => {
    if (session && !hasFetchedRef.current) {
      hasFetchedRef.current = true;
      fetchContacts();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [session]);

  const fetchContacts = async () => {
    if (!session) return;

    setLoading(true);
    try {
      const response = await fetch('/api/contacts');
      if (response.ok) {
        const data: Contact[] = await response.json();
        setContacts(data || []);
      }
    } finally {
      setLoading(false);
    }
  };

  const handleAddContact = async () => {
    if (!newContactCharacterName || !session) return;

    try {
      const response = await fetch('/api/contacts', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ characterName: newContactCharacterName }),
      });

      if (response.ok) {
        showSnackbar('Contact request sent!', 'success');
        setAddContactOpen(false);
        setNewContactCharacterName('');
        await fetchContacts();
      } else {
        const error = await response.json();
        showSnackbar(error.error || 'Failed to send contact request', 'error');
      }
    } catch (err) {
      showSnackbar('Failed to send contact request', 'error');
    }
  };

  const handleAccept = async (contactId: number) => {
    try {
      const response = await fetch(`/api/contacts/${contactId}/accept`, {
        method: 'POST',
      });

      if (response.ok) {
        showSnackbar('Contact accepted!', 'success');
        await fetchContacts();
      } else {
        const error = await response.json();
        showSnackbar(error.error || 'Failed to accept contact', 'error');
      }
    } catch (err) {
      showSnackbar('Failed to accept contact', 'error');
    }
  };

  const handleReject = async (contactId: number) => {
    try {
      const response = await fetch(`/api/contacts/${contactId}/reject`, {
        method: 'POST',
      });

      if (response.ok) {
        showSnackbar('Contact rejected', 'success');
        await fetchContacts();
      } else {
        const error = await response.json();
        showSnackbar(error.error || 'Failed to reject contact', 'error');
      }
    } catch (err) {
      showSnackbar('Failed to reject contact', 'error');
    }
  };

  const handleDelete = async (contactId: number) => {
    try {
      const response = await fetch(`/api/contacts/${contactId}`, {
        method: 'DELETE',
      });

      if (response.ok) {
        showSnackbar('Contact removed', 'success');
        await fetchContacts();
      } else {
        const error = await response.json();
        showSnackbar(error.error || 'Failed to remove contact', 'error');
      }
    } catch (err) {
      showSnackbar('Failed to remove contact', 'error');
    }
  };

  const handleOpenPermissions = (contact: Contact) => {
    setSelectedContact(contact);
    setPermissionsDialogOpen(true);
  };

  const handleClosePermissions = () => {
    setPermissionsDialogOpen(false);
    setSelectedContact(null);
  };

  const showSnackbar = (message: string, severity: 'success' | 'error') => {
    setSnackbarMessage(message);
    setSnackbarSeverity(severity);
    setSnackbarOpen(true);
  };

  // Filter contacts by tab
  const myContacts = contacts.filter(c => c.status === 'accepted');
  const pendingRequests = contacts.filter(c =>
    c.status === 'pending' && c.recipientUserId === currentUserId
  );
  const sentRequests = contacts.filter(c =>
    c.status === 'pending' && c.requesterUserId === currentUserId
  );

  if (!session) {
    return null;
  }

  if (loading) {
    return <Loading />;
  }

  return (
    <>
      <Navbar />
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h4">Contacts</Typography>
          <Button
            variant="contained"
            startIcon={<PersonAddIcon />}
            onClick={() => setAddContactOpen(true)}
          >
            Add Contact
          </Button>
        </Box>

        <Card>
          <Tabs value={tabIndex} onChange={(_, newValue) => setTabIndex(newValue)}>
            <Tab label={`My Contacts (${myContacts.length})`} />
            <Tab label={`Pending Requests (${pendingRequests.length})`} />
            <Tab label={`Sent Requests (${sentRequests.length})`} />
          </Tabs>

          <CardContent>
            {/* My Contacts Tab */}
            {tabIndex === 0 && (
              <TableContainer component={Paper} variant="outlined">
                {myContacts.length === 0 ? (
                  <Box sx={{ p: 4, textAlign: 'center' }}>
                    <Typography color="text.secondary">
                      No contacts yet. Add a contact to get started!
                    </Typography>
                  </Box>
                ) : (
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Character Name</TableCell>
                        <TableCell>Status</TableCell>
                        <TableCell>Connected Since</TableCell>
                        <TableCell align="right">Actions</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {myContacts.map((contact) => {
                        const otherUserName = contact.requesterUserId === currentUserId
                          ? contact.recipientName
                          : contact.requesterName;

                        return (
                          <TableRow key={contact.id} hover>
                            <TableCell>{otherUserName}</TableCell>
                            <TableCell>
                              <Chip label="Connected" color="success" size="small" />
                            </TableCell>
                            <TableCell>
                              {new Date(contact.respondedAt || contact.requestedAt).toLocaleDateString()}
                            </TableCell>
                            <TableCell align="right">
                              <IconButton
                                size="small"
                                onClick={() => handleOpenPermissions(contact)}
                                title="Manage Permissions"
                              >
                                <SettingsIcon />
                              </IconButton>
                              <IconButton
                                size="small"
                                onClick={() => handleDelete(contact.id)}
                                title="Remove Contact"
                                color="error"
                              >
                                <DeleteIcon />
                              </IconButton>
                            </TableCell>
                          </TableRow>
                        );
                      })}
                    </TableBody>
                  </Table>
                )}
              </TableContainer>
            )}

            {/* Pending Requests Tab */}
            {tabIndex === 1 && (
              <TableContainer component={Paper} variant="outlined">
                {pendingRequests.length === 0 ? (
                  <Box sx={{ p: 4, textAlign: 'center' }}>
                    <Typography color="text.secondary">
                      No pending requests
                    </Typography>
                  </Box>
                ) : (
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Character Name</TableCell>
                        <TableCell>Requested</TableCell>
                        <TableCell align="right">Actions</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {pendingRequests.map((contact) => (
                        <TableRow key={contact.id} hover>
                          <TableCell>{contact.requesterName}</TableCell>
                          <TableCell>
                            {new Date(contact.requestedAt).toLocaleDateString()}
                          </TableCell>
                          <TableCell align="right">
                            <IconButton
                              size="small"
                              onClick={() => handleAccept(contact.id)}
                              title="Accept"
                              color="success"
                            >
                              <CheckIcon />
                            </IconButton>
                            <IconButton
                              size="small"
                              onClick={() => handleReject(contact.id)}
                              title="Reject"
                              color="error"
                            >
                              <CloseIcon />
                            </IconButton>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}
              </TableContainer>
            )}

            {/* Sent Requests Tab */}
            {tabIndex === 2 && (
              <TableContainer component={Paper} variant="outlined">
                {sentRequests.length === 0 ? (
                  <Box sx={{ p: 4, textAlign: 'center' }}>
                    <Typography color="text.secondary">
                      No sent requests
                    </Typography>
                  </Box>
                ) : (
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Character Name</TableCell>
                        <TableCell>Sent</TableCell>
                        <TableCell align="right">Actions</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {sentRequests.map((contact) => (
                        <TableRow key={contact.id} hover>
                          <TableCell>{contact.recipientName}</TableCell>
                          <TableCell>
                            {new Date(contact.requestedAt).toLocaleDateString()}
                          </TableCell>
                          <TableCell align="right">
                            <IconButton
                              size="small"
                              onClick={() => handleDelete(contact.id)}
                              title="Cancel Request"
                              color="error"
                            >
                              <CloseIcon />
                            </IconButton>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}
              </TableContainer>
            )}
          </CardContent>
        </Card>
      </Container>

      {/* Add Contact Dialog */}
      <Dialog open={addContactOpen} onClose={() => setAddContactOpen(false)}>
        <DialogTitle>Add Contact</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Character Name"
            type="text"
            fullWidth
            value={newContactCharacterName}
            onChange={(e) => setNewContactCharacterName(e.target.value)}
            helperText="Enter the character name of the person you want to add"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setAddContactOpen(false)}>Cancel</Button>
          <Button onClick={handleAddContact} variant="contained">
            Send Request
          </Button>
        </DialogActions>
      </Dialog>

      {/* Permissions Dialog */}
      {selectedContact && (
        <PermissionsDialog
          open={permissionsDialogOpen}
          onClose={handleClosePermissions}
          contact={selectedContact}
          currentUserId={currentUserId || 0}
        />
      )}

      {/* Snackbar */}
      <Snackbar
        open={snackbarOpen}
        autoHideDuration={3000}
        onClose={() => setSnackbarOpen(false)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert
          onClose={() => setSnackbarOpen(false)}
          severity={snackbarSeverity}
          sx={{ width: '100%' }}
        >
          {snackbarMessage}
        </Alert>
      </Snackbar>
    </>
  );
}
