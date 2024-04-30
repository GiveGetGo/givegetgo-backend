import React, { useState, useEffect } from 'react';
import { StyleSheet, View, Text, TextInput, TouchableOpacity, SafeAreaView } from 'react-native';
import { Appbar } from 'react-native-paper';
import { MaterialCommunityIcons } from '@expo/vector-icons';
import { useNavigation } from '@react-navigation/native';
import { StackNavigationProp } from '@react-navigation/stack';
import { useFonts, Montserrat_700Bold_Italic } from '@expo-google-fonts/montserrat';
import requestMFASetup from './SignUpScreen'

type RootStackParamList = {
  CheckEmailScreen: undefined;
  ConfirmationScreen: undefined;
};

type ScreenNavigationProp = StackNavigationProp<
  RootStackParamList,
  'CheckEmailScreen' | 'ConfirmationScreen'
>;

const CheckEmailScreen: React.FC = () => {

  const [fontsLoaded] = useFonts({ Montserrat_700Bold_Italic });

  const use_navigation = useNavigation(); //for Appbar.BackAction

  const navigation = useNavigation<ScreenNavigationProp>();
  const [code, setCode] = useState('');
  const [email, setEmail] = useState<string>('xxx@gmail.com');

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
        verifyEmailCode(email, 'register', text)
      }
    }
  };

  const handleConfirm = () => {
    navigation.navigate('ConfirmationScreen');
  };

  const handleResendCode = () => {
    requestMFASetup
  };

  async function verifyEmailCode(email: string, event: string, verificationCode: string) {
    try {
      const response = await fetch('http://api.givegetgo.xyz/v1/verification/verify-email', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: email,
          event: event,
          verification_code: verificationCode
        }),
      });
  
      const json = await response.json(); // Parse the JSON response
      console.log("Verification response:", json);
  
      if (response.status === 200) {
        console.log('Email verified successfully:', json);
        alert('Email verified successfully!');
        handleConfirm();
        // You can navigate to another screen or update the UI accordingly
      } else if (response.status === 400) {
        console.error('Bad request:', json.msg);
        alert(`Error: ${json.msg}`);
      } else if (response.status === 401) {
        console.error('Unauthorized:', json.msg);
        alert(`Error: ${json.msg}`);
      } else if (response.status === 500) {
        console.error('Internal server error:', json.msg);
        alert(`Error: ${json.msg}`);
      } else {
        console.error('Unexpected error:', json);
        alert(`Error: ${json.msg}`);
      }
    } catch (error) {
      console.error('Network error:', error);
      alert('Failed to connect to the server. Please try again later.');
    }
  }

  useEffect(() => {
    // Fetch the email from the backend
    const fetchEmail = async () => {                        
      try {
        const response = await fetch('http://api.givegetgo.xyz/v1/user/me');
        const json = await response.json();
        setEmail(json.email); // Adjust this depending on the structure of your JSON
      } catch (error) {
        // console.error(error);
      }
    };

    fetchEmail();
  }, []);

return (
  <SafeAreaView style={styles.container}>
    <View style={styles.headerContainer}>
        <Appbar.BackAction style={styles.backAction} onPress={() => use_navigation.goBack()} />
        <Text style={styles.header}>GiveGetGo</Text>
        <View style={styles.backActionPlaceholder} />
    </View>
    <MaterialCommunityIcons name="email-outline" size={100} color="#000" />
    <Text style={styles.title}>Check Your Email</Text>
    <Text style={styles.subtitle1}>
      We have sent an email to {email} to confirm the validity of this email address.
    </Text>
    <Text style={styles.subtitle2}>
      Please enter the 7-digit code below.
    </Text>
    <TextInput
      style={styles.codeInput}
      placeholder="_ _ _ _ _ _ _"
      value={code}
      onChangeText={handleCodeChange}
      keyboardType="number-pad"
      maxLength={7}
      returnKeyType="done"
    />
    <TouchableOpacity style={styles.button} onPress={handleResendCode}>
      <Text style={styles.buttonText}>Resend Code</Text>
    </TouchableOpacity>
  </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,                                
    marginTop: 50,
    justifyContent: 'center',
    alignItems: 'center',
  },
  headerContainer: {
    flexDirection: 'row', // Aligns items in a row
    alignItems: 'center', // Centers items vertically
    justifyContent: 'space-between', // Distributes items evenly horizontally
    paddingLeft: 10, 
    paddingRight: 10, 
    position: 'absolute', // So that while setting card to the vertical middle, it still stays at the same place
    top: 0, 
    left: 0,
    right: 0,
    zIndex: 1, // Ensure the headerContainer is above the card
  },
  header: {
    fontSize: 22, // Increase the font size
    fontWeight: '600', // Make the font weight bold
    fontFamily: 'Montserrat_700Bold_Italic',
    textAlign: 'center', // Center the text
    color: '#444444', // Dark gray color
  },
  backActionPlaceholder: {
    width: 48, // This should match the width of the Appbar.BackAction for balance
    height: 48,
  },
  backAction: {
    marginLeft: 0 //This means the relative margin, comparing to the container (?)
  },
  emailIcon: {
    marginBottom: 24,
  },
  title: {
    fontSize: 22,
    fontWeight: 'bold',
    marginVertical: 8,
  },
  subtitle1: {
    fontSize: 16,
    color: 'grey',
    textAlign: 'center',
    // marginBottom: 24,
    padding: 20,
  },
  subtitle2: {
    fontSize: 16,
    color: 'grey',
    textAlign: 'center',
    marginTop: -40,
    marginBottom: 24,
    padding: 20,
  },
  codeInput: {
    width: '100%',
    padding: 10,
    fontSize: 18,
    borderBottomColor: 'grey',
    textAlign: 'center',
    marginBottom: 24,
  },
  button: {
    width: '80%',
    padding: 12,
    borderRadius: 5,
    backgroundColor: 'black',
    alignItems: 'center',
    marginBottom: 24,
  },
  buttonText: {
    fontSize: 18,
    color: 'white',
    fontWeight: '500',
  },
});

export default CheckEmailScreen;
