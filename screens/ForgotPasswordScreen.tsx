import React, { useState } from 'react';
import { StyleSheet, View, Text, TextInput, TouchableOpacity } from 'react-native';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';

// Define the types for your navigation stack
type RootStackParamList = {
  ForgotPasswordScreen: undefined;
  CheckEmailScreen: undefined;
  SignUpScreen: undefined;
};

// Define the type for the navigation prop
type ScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'ForgotPasswordScreen' | 'CheckEmailScreen' | 'SignUpScreen'
>;

const ForgotPasswordScreen: React.FC = () => {
  const navigation = useNavigation<ScreenNavigationProp>();
  const [email, setEmail] = useState<string>('');

  const handleResetPassword = () => {
    // Handle the reset password logic
    navigation.navigate('CheckEmailScreen');
    console.log('Reset password for:', email);
  };

  const handleSignUp = () => {
    navigation.navigate('SignUpScreen');
  };

  return (
    <View style={styles.container}>
      <Text style={styles.titleText}>Oh, No! I Forgot</Text>
      <Text style={styles.subtitleText}>Enter your email and we'll send you a link to change a new password</Text>

      <TextInput
        style={styles.input}
        placeholder="Email"
        onChangeText={setEmail}
        value={email}
        keyboardType="email-address"
        autoCapitalize="none"
      />

      <TouchableOpacity style={styles.button} onPress={handleResetPassword}>
        <Text style={styles.buttonText}>Forgot Password</Text>
      </TouchableOpacity>

      <Text style={styles.text}>
        Don't have an account?{' '}
        <TouchableOpacity onPress={handleSignUp}>
          <Text style={styles.linkText}>Sign Up</Text>
        </TouchableOpacity>
      </Text>
    </View>
  );
};

// You may want to customize these styles to match the UI provided in the Figma design
const styles = StyleSheet.create({
  // ... other styles remain unchanged
  titleText: {
    fontSize: 32,
    fontWeight: 'bold',
    marginBottom: 8,
  },
  subtitleText: {
    fontSize: 16,
    marginBottom: 48,
  },
});

export default ForgotPasswordScreen;