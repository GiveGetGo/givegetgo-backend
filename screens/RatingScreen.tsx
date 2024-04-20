import React, { useState } from 'react';
import { View, StyleSheet, TouchableOpacity, SafeAreaView } from 'react-native';
import { Avatar, Button, Text, Card, Title, Paragraph, Appbar } from 'react-native-paper';
import { Rating } from 'react-native-ratings';
import { useNavigation } from '@react-navigation/native';
import { NativeStackScreenProps } from '@react-navigation/native-stack';

type RootStackParamList = {
  NotificationScreen: undefined;
  RatingScreen: undefined;
  RatingSucceedScreen: undefined;
};

type NotificationsProps = NativeStackScreenProps<RootStackParamList, 'NotificationScreen'>;

const RatingScreen: React.FC<NotificationsProps> = ({ navigation }: NotificationsProps) => {
  const use_navigation = useNavigation(); //for Appbar.BackAction

  const [rating, setRating] = useState<number>(0);

  const submitRating = () => {
    // Submit the rating to your backend or handle 
    console.log('Rating submitted: ', rating);
    navigation.navigate('RatingSucceedScreen');
  };

  return (
    <SafeAreaView  style={styles.container}> 
      <View style={styles.headerContainer}>
        <Appbar.BackAction style={styles.backAction} onPress={() => use_navigation.goBack()} />
        <Text style={styles.header}>GiveGetGo</Text>
        <View style={styles.backActionPlaceholder} />
      </View>
      <View style={styles.rating_container}>
        <Text style={styles.question}>
          How would you rate your match with <Paragraph style={styles.boldQuestion}>xxxxxx</Paragraph>?
        </Text>
        <Rating
          type="custom"
          ratingCount={5}
          imageSize={40} // Smaller size for the stars
          fractions={0} // Set the granularity of the rating; 1 means full stars
          startingValue={rating} // The initial rating value
          tintColor="#f6f6f6" // Background color for the rating component
          ratingBackgroundColor="#c8c7c8" // Color behind the non-selected rating icons
          ratingColor="orange" // Color of the selected rating icons
          showRating
          onFinishRating={(rating: number) => setRating(rating)}
          style={styles.rating}
        />
        <Button style={styles.button} mode="contained" onPress={submitRating}>
          Submit
        </Button>
      </View>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    marginTop: 50,
    alignItems: 'center',
  },
  header: {
    fontSize: 20,
    fontWeight: 'bold',
    padding: 16,
    alignItems: 'center',
  },
  headerContainer: {
    flexDirection: 'row', // Aligns items in a row
    alignItems: 'center', // Centers items vertically
    paddingLeft: 10, // Adds padding to the left of the avatar
    paddingRight: 10, // Adds padding to the right side
  },
  backActionPlaceholder: {
    width: 48, 
    height: 48,
  },
  backAction: {
    marginLeft: 0 
  },
  button: {
    textAlign: 'center',
    marginBottom: 50,
    marginTop: 20,
  //   position: 'absolute', 
  //   left: 80,
  //   right: 80, //position, left, right together controls the button's length and horizontal location
  //   alignSelf: 'center', 
  },
  question: {
    fontSize: 18,
    textAlign: 'center',
    marginBottom: 8,
  },
  boldQuestion: {
    fontSize: 18,
    textAlign: 'center',
    marginBottom: 8,
    fontWeight: 'bold',
  },
  rating: {
    paddingVertical: 10,
  },
  rating_container: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
    padding: 24,
    marginTop: -20, // the height of all three components
  },
});

export default RatingScreen;
