import React from 'react';
import { StyleSheet, SafeAreaView, FlatList, View} from 'react-native';
import { Text, Avatar, Divider, Card, Title, Paragraph, Appbar} from 'react-native-paper';
import { Rating } from 'react-native-ratings';
import { useNavigation } from '@react-navigation/native';

// Define the types for your navigation stack
type RootStackParamList = {
  SignUpScreen: undefined;
  CheckEmailScreen: undefined;
  LoginScreen: undefined;
};

type UserInfo = {
  name: string;
  email: string;
  bio: string;
  profilePicture: string; // This should be a URI string
  major: string;
  classYear: number;
};

type History = {
  id: string;
  title: string;
  subtitle: string;
  description: string;
  time: string;
  rating?: number; // Optional rating, only for "Match" items
};

const currentUserInfo: UserInfo = { //should be loaded from the database
  name: 'Jimmy Ho',
  email: 'xxx@gmail.com',
  bio: 'Hi I am Jimmy.',
  profilePicture: '', 
  major: 'Computer Engineering',
  classYear: 2024,
};

const history_data: History[] = [ //should be loaded from the database //[] because there are more than one data
  {
    id: '1',
    title: 'Match',
    subtitle: '$10 Panera Giftcard ',
    description: 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx...',
    time: '1 hour ago',
    rating: 3,
  },
  {
    id: '2',
    title: 'Request',
    subtitle: '$10 Panera Giftcard ',
    description: 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx...',
    time: '2 hour ago',
  },
];

const NotificationStackProfileScreen: React.FC = () => {

  const use_navigation = useNavigation(); //for Appbar.BackAction

  const renderItem = ({ item }: { item: History })  => (
    <Card style={styles.history_card}>
      <Card.Content>
        <View style={styles.titleRow}>
          <Title style={styles.history_title}>{item.title}</Title>
          <Paragraph style={styles.history_time}>{item.time}</Paragraph>
        </View>
        <View style={styles.subtitleRow}>
          <Paragraph style={styles.history_subtitle}>{item.subtitle}</Paragraph>
          {item.title === 'Match' && item.rating !== undefined && (
            <Rating
              type='custom' //so that "ratingBackgroundColor" works
              ratingCount={5}
              imageSize={20}
              readonly
              startingValue={item.rating}
              tintColor="#f6f6f6"
              ratingBackgroundColor='#c8c7c8' 
            />
          )}
        </View>
        <Paragraph style={styles.history_description}>{item.description}</Paragraph>
      </Card.Content>
    </Card>
  );

  return (
    <SafeAreaView  style={styles.container}>
      <View style={styles.headerContainer}>
        <Appbar.BackAction style={styles.backAction} onPress={() => use_navigation.goBack()} />
        <Text style={styles.header}>GiveGetGo</Text>
        <View style={styles.backActionPlaceholder} />
      </View>
      <Avatar.Image source={require('./profile_icon.jpg')} style={styles.settings_avatar} /> 
      <Text style={styles.name}>{currentUserInfo.name}</Text>
      <Text style={styles.details}>{`Class of ${currentUserInfo.classYear} • ${currentUserInfo.major}`}</Text>
      <Text style={styles.details_bio}>{currentUserInfo.bio}</Text>
      <Divider style={styles.divider} />
      <FlatList
          data={history_data}
          renderItem={renderItem}
          keyExtractor={(item) => item.id}
          // removeClippedSubviews={true} // If you want to improve performance on large lists, uncomment the line
      />
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    marginTop: 30,
  },
  header: {
    fontSize: 20,
    fontWeight: 'bold',
    padding: 16,
  },
  settings_avatar: {
    alignSelf: 'center',
    marginTop: 20,
    backgroundColor: '#c7c7c7', // Placeholder color
  },
  name: {
      fontSize: 20,
      fontWeight: 'bold',
      textAlign: 'center',
      marginTop: 8,
  },
  details: {
      textAlign: 'center',
      marginBottom: 6,
  },
  details_bio: {
      textAlign: 'center',
      fontStyle: 'italic',
      marginBottom: 6,
  },
  divider: {
      marginVertical: 8,
  },
  history_card: { //page gets longer when there are more contexts
    borderRadius: 15, // Add rounded corners to the card
    marginVertical: 6,
    marginHorizontal: 10,
    // elevation: 1, // Adjust for desired shadow depth
    // backgroundColor: '#ffffff', 
    padding: 6, // Add padding inside the card
  },
  titleRow: { //binding title and time together on the same row
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
  subtitleRow: { //binding subtitle and rating together on the same row
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 4, 
  },
  history_title: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  history_subtitle: {
    fontSize: 14, // Adjust the font size as needed
    color: '#000', // Adjust the color as needed
  },
  history_time: {
    fontSize: 12,
    color: 'gray', // Adjust the color for the time text
  },
  history_description: {
    fontSize: 14,
    marginBottom: 4, // Adjust spacing to match the image provided
  },
  backActionPlaceholder: {
    width: 48, // This should match the width of the Appbar.BackAction for balance
    height: 48,
  },
  backAction: {
    marginLeft: 0 //This means the relative margin, comparing to the container (?)
  },
  headerContainer: {
    flexDirection: 'row', // Aligns items in a row
    alignItems: 'center', // Centers items vertically
    paddingLeft: 10, // Adds padding to the left of the avatar
    paddingRight: 10, // Adds padding to the right side
  },
});

export default NotificationStackProfileScreen;
