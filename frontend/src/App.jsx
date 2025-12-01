import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './hooks/useAuth';
import Login from './components/Auth/Login';
import Register from './components/Auth/Register';
import Dashboard from './components/Dashboard/Dashboard';
import Profile from './components/Profile/Profile';
import Groups from './components/Groups/Groups';
import GroupInvitations from './components/Groups/GroupInvitations';
import GroupSettings from './components/Groups/GroupSettings';
import TaskDetail from './components/Task/TaskDetail';
import CreateGroupTask from './components/Task/CreateGroupTask';
import ProtectedRoute from './components/ProtectedRoute';
import './App.css';

function App() {
    const { token, loading } = useAuth();

    if (loading) {
        return (
            <div style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                height: '100vh',
                fontSize: '18px',
                color: '#6b7280'
            }}>
                🔄 بارگذاری...
            </div>
        );
    }

    return (
        <div className="app-container">
            <Routes>
                {/* Public routes */}
                <Route
                    path="/login"
                    element={token ? <Navigate to="/dashboard" replace /> : <Login />}
                />
                <Route
                    path="/register"
                    element={token ? <Navigate to="/dashboard" replace /> : <Register />}
                />

                {/* Protected routes */}
                <Route path="/dashboard" element={
                    <ProtectedRoute>
                        <Dashboard />
                    </ProtectedRoute>
                } />

                <Route path="/profile" element={
                    <ProtectedRoute>
                        <Profile />
                    </ProtectedRoute>
                } />

                <Route path="/groups" element={
                    <ProtectedRoute>
                        <Groups />
                    </ProtectedRoute>
                } />

                <Route path="/groups/invitations" element={
                    <ProtectedRoute>
                        <GroupInvitations />
                    </ProtectedRoute>
                } />

                <Route path="/groups/:id/settings" element={
                    <ProtectedRoute>
                        <GroupSettings />
                    </ProtectedRoute>
                } />

                <Route path="/groups/:id/add-task" element={
                    <ProtectedRoute>
                        <CreateGroupTask />
                    </ProtectedRoute>
                } />

                <Route path="/tasks/:id" element={
                    <ProtectedRoute>
                        <TaskDetail />
                    </ProtectedRoute>
                } />

                {/* Default routes */}
                <Route
                    path="/"
                    element={<Navigate to={token ? "/dashboard" : "/login"} replace />}
                />

                {/* Catch all route */}
                <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
        </div>
    );
}

export default App;