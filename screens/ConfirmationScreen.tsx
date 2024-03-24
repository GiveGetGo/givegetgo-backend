import React, { useState, useEffect } from 'react';
import { View, StyleSheet, Image, Text, TouchableOpacity } from 'react-native';
import { Button } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';

// Define the types for your navigation stack
type RootStackParamList = {
  ConfirmationScreen: undefined;
  LoginScreen: undefined;
};

// Define the type for the navigation prop
type LoginScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'ConfirmationScreen' | 'LoginScreen'
>;

const ConfirmationScreen: React.FC = () => {
  const checkIcon = require('./confirm_icon.jpg'); 
  const navigation = useNavigation<LoginScreenNavigationProp>();
  const [email, setEmail] = useState<string>('');

  const handleGoHome = () => {
    // Navigation logic to go back to the home screen
    navigation.navigate('LoginScreen');
  };

  useEffect(() => {
    // Fetch the email from the backend
    const fetchEmail = async () => {
      try {
        const response = await fetch('URL_TO_YOUR_BACKEND/json_endpoint');
        const json = await response.json();
        setEmail(json.email); // Adjust this depending on the structure of your JSON
      } catch (error) {
        console.error(error);
      }
    };

    fetchEmail();
  }, []);

  return (
    <View style={styles.container}>
      <Text style={styles.header}>GiveGetGo</Text>
      <Image source={checkIcon} style={styles.icon} />
      <Text style={styles.confirmedText}>Confirmed</Text>
      <Text style={styles.emailText}>{email} has been confirmed</Text> 
      <Button mode="contained" onPress={handleGoHome} style={styles.button}>
        Home
      </Button>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 20,
    backgroundColor: '#fff', // adjust the background color as per your theme
  },
  header: {
    position: 'absolute',
    top: 20, // adjust the positioning as needed
    left: 20, // adjust the positioning as needed
    fontSize: 20,
    fontWeight: 'bold',
  },
  icon: {
    width: 100, // Set the width as per your UI design
    height: 100, // Set the height as per your UI design
    marginBottom: 24, // adjust the spacing as per your UI design
  },
  confirmedText: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 8,
  },
  emailText: {
    fontSize: 16,
    marginBottom: 48,
  },
  button: {
    // Style your button with react-native-paper theming or custom styles
  },
});

export default ConfirmationScreen;
