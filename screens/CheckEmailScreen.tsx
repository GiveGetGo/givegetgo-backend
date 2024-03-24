import React, { useState, useEffect } from 'react';
import { StyleSheet, View, Text, TextInput, Image, TouchableOpacity } from 'react-native';
import { Button } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';

type RootStackParamList = {
  // ... other screen names
  CheckEmailScreen: undefined;
  ConfirmationScreen: undefined;
};

type ScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'CheckEmailScreen' | 'ConfirmationScreen'
>;

const CheckEmailScreen: React.FC = () => {
  const navigation = useNavigation<ScreenNavigationProp>();
  const [code, setCode] = useState('');
  const [email, setEmail] = useState<string>('');

  const handleCodeComplete = (code: string) => {
    // This function is triggered when all seven digits are entered
    console.log('Code entered:', code);
    // Add your logic here, for example navigate to a new screen or verify the code
  };
  const handleCodeChange = (text: string) => {
    if (text.length <= 7) {
      setCode(text);
      if (text.length === 7) {
        handleCodeComplete(text);
        handleConfirm();
      }
    }
  };

  const handleConfirm = () => {
    navigation.navigate('ConfirmationScreen');
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
      <TouchableOpacity onPress={() => navigation.goBack()} style={styles.backButton}>
        <Text>←</Text> 
      </TouchableOpacity>
      <Text style={styles.header}>GiveGetGo</Text>
      <Image
        source={require('./email_icon.png')} // Replace with your email icon path
        style={styles.emailIcon}
      />
      <Text style={styles.title}>Check Your Email</Text>
      <Text style={styles.subtitle}>
        We have sent an email to {email} to confirm the validity of this email address.
        After receiving the email, please enter the 7-digit code in the provided box.
      </Text>
      <View style={styles.container}>
        <TextInput
            style={styles.codeInput}
            placeholder="_ _ _ _ _ _ _"
            value={code}
            onChangeText={handleCodeChange}
            keyboardType="number-pad"
            maxLength={7}
            returnKeyType="done"
        />
      </View>
      <Button
        mode="outlined"
        onPress={() => {
          // Handle resend code logic
        }}
      >
        Resend Code
      </Button>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: 'center',
    padding: 20,
    backgroundColor: '#fff',
  },
  backButton: {
    alignSelf: 'flex-start',
    marginBottom: 20,
  },
  header: {
    // Style for your header
  },
  emailIcon: {
    width: 100,
    height: 100,
    // Adjust the size as needed
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    marginVertical: 16,
  },
  subtitle: {
    textAlign: 'center',
    marginBottom: 40,
  },
});

export default CheckEmailScreen;