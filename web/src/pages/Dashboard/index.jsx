/*
Copyright (C) 2025 Auth Gateway
*/

import React, { useContext } from 'react';
import { Card, Typography, Avatar, Descriptions } from '@douyinfe/semi-ui';
import { IconUser, IconShield, IconKey } from '@douyinfe/semi-icons';
import { UserContext } from '../../context/User';
import { Link } from 'react-router-dom';

const { Title, Text } = Typography;

const Dashboard = () => {
  const [userState] = useContext(UserContext);
  const user = userState?.user || {};

  return (
    <div className='mt-[60px] px-4 max-w-4xl mx-auto'>
      <Title heading={2} className='mb-6'>Dashboard</Title>

      <div className='grid grid-cols-1 md:grid-cols-2 gap-4'>
        {/* User Info Card */}
        <Card title='Profile' className='h-fit'>
          <div className='flex items-center gap-4 mb-4'>
            <Avatar size='large' color='blue'>
              {user.username?.[0]?.toUpperCase() || 'U'}
            </Avatar>
            <div>
              <Text strong>{user.display_name || user.username}</Text>
              <br />
              <Text type='tertiary'>{user.email || 'No email set'}</Text>
            </div>
          </div>
          <Descriptions
            data={[
              { key: 'Username', value: user.username },
              { key: 'Role', value: user.role === 100 ? 'Admin' : user.role === 10 ? 'User' : 'Guest' },
              { key: 'Status', value: user.status === 1 ? 'Active' : 'Disabled' },
            ]}
          />
          <div className='mt-4'>
            <Link to='/console/personal' className='text-blue-500 hover:underline'>
              Edit Profile
            </Link>
          </div>
        </Card>

        {/* Security Card */}
        <Card title='Security' className='h-fit'>
          <div className='space-y-3'>
            <div className='flex items-center gap-2'>
              <IconKey />
              <Text>Passkey: {user.passkey_enabled ? 'Enabled' : 'Not set'}</Text>
            </div>
            <div className='flex items-center gap-2'>
              <IconShield />
              <Text>2FA: {user.twofa_enabled ? 'Enabled' : 'Not set'}</Text>
            </div>
          </div>
          <div className='mt-4'>
            <Link to='/console/personal' className='text-blue-500 hover:underline'>
              Security Settings
            </Link>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default Dashboard;
